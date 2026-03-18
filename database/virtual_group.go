package database

import (
	"context"

	"github.com/forbole/bdjuno/v4/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) SaveGVG(ctx context.Context, gvg *models.GlobalVirtualGroup) error {
	return nil
}

func (db *DB) UpdateGVG(ctx context.Context, gvg *models.GlobalVirtualGroup) error {
	return nil
}

func (db *DB) SaveLVG(ctx context.Context, lvg *models.LocalVirtualGroup) error {
	return nil
}

func (db *DB) UpdateLVG(ctx context.Context, lvg *models.LocalVirtualGroup) error {
	return nil
}

func (db *DB) SaveVGF(ctx context.Context, vgf *models.GlobalVirtualGroupFamily) error {
	return nil
}

func (db *DB) UpdateVGF(ctx context.Context, vgf *models.GlobalVirtualGroupFamily) error {
	return nil
}

func (db *DB) SaveGVGToSQL(ctx context.Context, gvg *models.GlobalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.GlobalVirtualGroup{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "global_virtual_group_id"}},
		UpdateAll: true,
	}).Create(gvg).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateGVGToSQL(ctx context.Context, gvg *models.GlobalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Model(&models.GlobalVirtualGroup{}).
		Select("primary_sp_id", "secondary_sp_ids", "stored_size", "total_deposit", "update_at", "update_tx_hash", "update_time").
		Where("global_virtual_group_id = ?", gvg.GlobalVirtualGroupID).
		Updates(gvg).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) DeleteGVGToSQL(ctx context.Context, gvg *models.GlobalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Model(&models.GlobalVirtualGroup{}).
		Select("removed", "update_at", "update_tx_hash", "update_time").
		Where("global_virtual_group_id = ?", gvg.GlobalVirtualGroupID).
		Updates(gvg).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) SaveLVGToSQL(ctx context.Context, lvg *models.LocalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.LocalVirtualGroup{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bucket_id"}, {Name: "local_virtual_group_id"}},
		UpdateAll: true,
	}).Create(lvg).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) UpdateLVGToSQL(ctx context.Context, lvg *models.LocalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Model(&models.LocalVirtualGroup{}).
		Select("global_virtual_group_id", "stored_size", "update_at", "update_tx_hash", "update_time").
		Where("local_virtual_group_id = ? and bucket_id = ?", lvg.LocalVirtualGroupID, lvg.BucketID).
		Updates(lvg).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) DeleteLVGToSQL(ctx context.Context, lvg *models.LocalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Model(&models.LocalVirtualGroup{}).
		Select("removed", "update_at", "update_tx_hash", "update_time").
		Where("local_virtual_group_id = ? and bucket_id = ?", lvg.LocalVirtualGroupID, lvg.BucketID).
		Updates(lvg).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) SaveGVGFToSQL(ctx context.Context, vgf *models.GlobalVirtualGroupFamily) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.GlobalVirtualGroupFamily{}).TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "global_virtual_group_family_id"}},
		UpdateAll: true,
	}).Create(vgf).Statement
	return stat.SQL.String(), stat.Vars
}

// SyncGVGToSQL generates an upsert SQL that only updates fields available from chain query,
// preserving event-sourced fields (update_tx_hash, update_evm_tx_hash, update_time).
func (db *DB) SyncGVGToSQL(ctx context.Context, gvg *models.GlobalVirtualGroup) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.GlobalVirtualGroup{}).TableName()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "global_virtual_group_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"family_id", "primary_sp_id", "secondary_sp_ids",
			"stored_size", "virtual_payment_address", "total_deposit",
			"update_at", "removed",
		}),
	}).Create(gvg).Statement
	return stat.SQL.String(), stat.Vars
}

// SyncGVGFToSQL generates an upsert SQL that only updates fields available from chain query,
// preserving event-sourced fields (update_tx_hash, update_evm_tx_hash, update_time).
func (db *DB) SyncGVGFToSQL(ctx context.Context, vgf *models.GlobalVirtualGroupFamily) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.GlobalVirtualGroupFamily{}).TableName()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "global_virtual_group_family_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"primary_sp_id", "global_virtual_group_ids", "virtual_payment_address",
			"update_at", "removed",
		}),
	}).Create(vgf).Statement
	return stat.SQL.String(), stat.Vars
}

func (db *DB) GetGVGFamilyByID(ctx context.Context, familyID uint32) (*models.GlobalVirtualGroupFamily, error) {
	var vgf models.GlobalVirtualGroupFamily
	if err := db.G.WithContext(ctx).Where("global_virtual_group_family_id = ?", familyID).First(&vgf).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &vgf, nil
}

func (db *DB) GetGVGByID(ctx context.Context, gvgID uint32) (*models.GlobalVirtualGroup, error) {
	var gvg models.GlobalVirtualGroup
	if err := db.G.WithContext(ctx).Where("global_virtual_group_id = ?", gvgID).First(&gvg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &gvg, nil
}

func (db *DB) UpdateVGFToSQL(ctx context.Context, vgf *models.GlobalVirtualGroupFamily) (string, []interface{}) {
	stat := db.G.Session(&gorm.Session{DryRun: true}).Table((&models.GlobalVirtualGroupFamily{}).TableName()).Where("global_virtual_group_family_id = ?", vgf.GlobalVirtualGroupFamilyID).Updates(vgf).Statement
	return stat.SQL.String(), stat.Vars
}
