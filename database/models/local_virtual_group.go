package models

import (
	"time"
)

type LocalVirtualGroup struct {
	ID                   uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	LocalVirtualGroupID  uint32    `gorm:"column:local_virtual_group_id;index:idx_lvg_id;index:idx_lvg_bucket,priority:1"`
	GlobalVirtualGroupID uint32    `gorm:"column:global_virtual_group_id;index:idx_gvg_id"`
	BucketID             string    `gorm:"column:bucket_id;type:TEXT;index:idx_bucket_id;index:idx_lvg_bucket,priority:2"`
	StoredSize           uint64    `gorm:"column:stored_size"`
	CreateAt             int64     `gorm:"column:create_at"`
	CreateTxHash         string    `gorm:"column:create_tx_hash;type:TEXT;not null"`
	CreateEVMTxHash      string    `gorm:"column:create_evm_tx_hash;type:TEXT;not null"`
	CreateTime           time.Time `gorm:"column:create_time"`
	UpdateAt             int64     `gorm:"column:update_at"`
	UpdateTxHash         string    `gorm:"column:update_tx_hash;type:TEXT;not null"`
	UpdateEVMTxHash      string    `gorm:"column:update_evm_tx_hash;type:TEXT;not null"`
	UpdateTime           time.Time `gorm:"column:update_time"`
	Removed              bool      `gorm:"column:removed;default:false"`
}

func (*LocalVirtualGroup) TableName() string {
	return "local_virtual_groups"
}
