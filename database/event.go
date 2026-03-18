package database

import (
	"context"

	"github.com/forbole/bdjuno/v4/database/models"
	"gorm.io/gorm"
)

func (db *DB) SaveBucketEventToSQL(ctx context.Context, e *models.BucketEvent) (string, []interface{}) {
	e.Event = ExtractEvent(e.Event)
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.BucketEvent{}).TableName()).Create(e).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) SaveObjectEventToSQL(ctx context.Context, e *models.ObjectEvent) (string, []interface{}) {
	e.Event = ExtractEvent(e.Event)
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.ObjectEvent{}).TableName()).Create(e).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) SaveGroupEventToSQL(ctx context.Context, e *models.GroupEvent) (string, []interface{}) {
	e.Event = ExtractEvent(e.Event)
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.GroupEvent{}).TableName()).Create(e).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) SaveSPEventToSQL(ctx context.Context, e *models.SpEvent) (string, []interface{}) {
	e.Event = ExtractEvent(e.Event)
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.SpEvent{}).TableName()).Create(e).Statement
	return stat.SQL.String(), stat.Vars
}
