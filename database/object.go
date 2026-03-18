package database

import (
	"context"

	"github.com/forbole/bdjuno/v4/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) SaveObject(ctx context.Context, object *models.Object) error {
	return nil
}

func (db *DB) UpdateObject(ctx context.Context, object *models.Object) error {
	return nil
}

func (db *DB) GetObject(ctx context.Context, objectId string) (*models.Object, error) {
	var object models.Object
	if err := db.G.WithContext(ctx).Table((&models.Object{}).TableName()).Where(
		"object_id = ? AND removed IS NOT TRUE", objectId).First(&object).Error; err != nil {
		return nil, err
	}
	return &object, nil
}

func (db *DB) SaveObjectToSQL(ctx context.Context, object *models.Object) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Object{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "object_id"}},
		UpdateAll: true,
	}).Create(object).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateObjectToSQL(ctx context.Context, object *models.Object) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Object{}).TableName()).Where("object_id = ?", object.ObjectID).Updates(object).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateObjectByBucketNameAndObjectNameToSQL(ctx context.Context, object *models.Object) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.Object{}).TableName()).Where("bucket_name = ? AND object_name =? AND removed=false", object.BucketName, object.ObjectName).Updates(object).Statement
	return stat.SQL.String(), stat.Vars
}
