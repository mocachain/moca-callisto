package database

import (
	"context"
	"fmt"

	"github.com/forbole/bdjuno/v4/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) SaveBucketToSQL(ctx context.Context, bucket *models.Bucket) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Bucket{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bucket_id"}},
		UpdateAll: true,
	}).Create(bucket).Statement

	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateBucketToSQL(ctx context.Context, bucket *models.Bucket) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Bucket{}).TableName()).Where("bucket_id = ?", bucket.BucketID).Updates(bucket).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateBucketByNameToSQL(ctx context.Context, bucket *models.Bucket) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Bucket{}).TableName()).Where("bucket_name= ? AND removed=false", bucket.BucketName).Updates(bucket).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateBucketOffChainStatus(ctx context.Context, bucket *models.Bucket) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Bucket{}).TableName()).
		Where("bucket_name = ?", bucket.BucketName).
		Select("off_chain_status", "update_at", "update_tx_hash", "update_time").
		Updates(bucket).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateBucketSPAsDelegatedAgentToSQL(ctx context.Context, bucket *models.Bucket) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Bucket{}).TableName()).
		Where("bucket_id = ?", bucket.BucketID).
		Select("sp_as_delegated_agent_disabled", "update_at", "update_tx_hash", "update_evm_tx_hash", "update_time").
		Updates(bucket).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) BatchUpdateBucketSize(ctx context.Context, buckets []*models.Bucket) error {
	return db.G.Table((&models.Bucket{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bucket_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"storage_size", "charge_size"}),
	}).Create(buckets).Error
}

func (db *DB) UpdateStorageSizeToSQL(ctx context.Context, objectID, bucketName, operation string) (string, []interface{}) {
	tableName := (&models.Object{}).TableName()
	sql := `UPDATE buckets SET storage_size = storage_size %s CAST((SELECT payload_size FROM %s WHERE object_id = ?) AS DECIMAL(65,0)) WHERE bucket_name = ?`
	vars := []interface{}{objectID, bucketName}
	finalSQL := fmt.Sprintf(sql, operation, tableName)
	return finalSQL, vars
}

func (db *DB) UpdateChargeSizeToSQL(ctx context.Context, objectID, bucketName, operation string) (string, []interface{}) {
	tableName := (&models.Object{}).TableName()
	sql := `UPDATE buckets SET charge_size = charge_size %s CASE WHEN (CAST((SELECT payload_size FROM %s WHERE object_id = ?)AS DECIMAL(65,0)) < 128000) THEN CAST(128000 AS DECIMAL(65,0)) ELSE CAST((SELECT payload_size FROM %s WHERE object_id = ?) AS DECIMAL(65,0)) END WHERE bucket_name = ?`
	vars := []interface{}{objectID, objectID, bucketName}
	finalSql := fmt.Sprintf(sql, operation, tableName, tableName)
	return finalSql, vars
}

// GetBucketByBucketName get bucket by bucket name
func (db *DB) GetBucketByBucketName(ctx context.Context, bucketName string) (*models.Bucket, error) {
	var (
		bucket *models.Bucket
		err    error
	)
	err = db.G.WithContext(ctx).Table((&models.Bucket{}).TableName()).
		Where("bucket_name = ? AND removed = ?", bucketName, false).Take(&bucket).Error
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
