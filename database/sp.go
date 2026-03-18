package database

import (
	"context"

	"github.com/forbole/bdjuno/v4/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) CreateStorageProvider(ctx context.Context, storageProvider *models.StorageProvider) error {
	return nil
}

func (db *DB) UpdateStorageProvider(ctx context.Context, storageProvider *models.StorageProvider) error {
	return nil
}

func (db *DB) CreateStorageProviderToSQL(ctx context.Context, storageProvider *models.StorageProvider) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.StorageProvider{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sp_id"}},
		UpdateAll: true,
	}).Create(storageProvider).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateStorageProviderToSQL(ctx context.Context, storageProvider *models.StorageProvider) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.StorageProvider{}).TableName()).Where("sp_id = ? ", storageProvider.SpID).Updates(storageProvider).Statement
	return stat.SQL.String(), stat.Vars
}

// SyncStorageProviderToSQL generates an upsert SQL that only updates fields available from chain query,
// preserving event-sourced fields (read_price, store_price, free_read_quota, update_time_sec, update_tx_hash, update_evm_tx_hash).
func (db *DB) SyncStorageProviderToSQL(ctx context.Context, storageProvider *models.StorageProvider) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.StorageProvider{}).TableName()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sp_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"operator_address", "funding_address", "seal_address", "approval_address", "gc_address",
			"total_deposit", "status", "endpoint",
			"moniker", "identity", "website", "security_contact", "details",
			"bls_key", "update_at", "removed",
		}),
	}).Create(storageProvider).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) GetStorageProviderBySpID(ctx context.Context, spID uint32) (*models.StorageProvider, error) {
	var sp models.StorageProvider
	err := db.G.Table((&models.StorageProvider{}).TableName()).Where("sp_id = ?", spID).First(&sp).Error
	if err != nil {
		return nil, err
	}
	return &sp, nil
}
