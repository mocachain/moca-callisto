package models

import (
	"time"

	"github.com/lib/pq"
)

type GlobalVirtualGroup struct {
	ID                    uint64        `gorm:"column:id;primaryKey;autoIncrement"`
	GlobalVirtualGroupID  uint32        `gorm:"column:global_virtual_group_id;index:idx_gvg_id"`
	FamilyID              uint32        `gorm:"column:family_id"`
	PrimarySpID           uint32        `gorm:"column:primary_sp_id;index:idx_primary_sp_id"`
	SecondarySpIDs        pq.Int32Array `gorm:"column:secondary_sp_ids;type:TEXT"`
	StoredSize            uint64        `gorm:"column:stored_size"`
	VirtualPaymentAddress string        `gorm:"column:virtual_payment_address;type:TEXT"`
	TotalDeposit          string        `gorm:"column:total_deposit;type:TEXT"`
	CreateAt              int64         `gorm:"column:create_at"`
	CreateTxHash          string        `gorm:"column:create_tx_hash;type:TEXT;not null"`
	CreateEVMTxHash       string        `gorm:"column:create_evm_tx_hash;type:TEXT;not null"`
	CreateTime            time.Time     `gorm:"column:create_time"`
	UpdateAt              int64         `gorm:"column:update_at"`
	UpdateTxHash          string        `gorm:"column:update_tx_hash;type:TEXT;not null"`
	UpdateEVMTxHash       string        `gorm:"column:update_evm_tx_hash;type:TEXT;not null"`
	UpdateTime            time.Time     `gorm:"column:update_time"`
	Removed               bool          `gorm:"column:removed;default:false"`
}

func (*GlobalVirtualGroup) TableName() string {
	return "global_virtual_groups"
}
