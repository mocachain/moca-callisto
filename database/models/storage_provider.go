package models

import (
	"fmt"
)

type StorageProvider struct {
	ID              uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	SpID            uint32 `gorm:"column:sp_id;index:idx_sp_id"`
	OperatorAddress string `gorm:"column:operator_address;type:TEXT;index:idx_operator_address"`
	FundingAddress  string `gorm:"column:funding_address;type:TEXT"`
	SealAddress     string `gorm:"column:seal_address;type:TEXT"`
	ApprovalAddress string `gorm:"column:approval_address;type:TEXT"`
	GcAddress       string `gorm:"column:gc_address;type:TEXT"`
	TotalDeposit    string `gorm:"column:total_deposit;type:TEXT"`
	Status          string `gorm:"column:status;type:VARCHAR(50)"`
	Endpoint        string `gorm:"column:endpoint;type:VARCHAR(256)"`
	Moniker         string `gorm:"column:moniker;type:VARCHAR(128)"`
	Identity        string `gorm:"column:identity;type:VARCHAR(256)"`
	Website         string `gorm:"column:website;type:VARCHAR(128)"`
	SecurityContact string `gorm:"column:security_contact;type:VARCHAR(128)"`
	Details         string `gorm:"column:details;type:VARCHAR(256)"`
	BlsKey          string `gorm:"column:bls_key;type:VARCHAR(96)"`
	UpdateTimeSec   int64  `gorm:"column:update_time_sec"`
	ReadPrice       string `gorm:"column:read_price;type:TEXT"`
	FreeReadQuota   uint64 `gorm:"column:free_read_quota"`
	StorePrice      string `gorm:"column:store_price;type:TEXT"`
	CreateAt        int64  `gorm:"column:create_at"`
	CreateTxHash    string `gorm:"column:create_tx_hash;type:TEXT;not null"`
	CreateEVMTxHash string `gorm:"column:create_evm_tx_hash;type:TEXT;not null"`
	UpdateAt        int64  `gorm:"column:update_at"`
	UpdateTxHash    string `gorm:"column:update_tx_hash;type:TEXT;not null"`
	UpdateEVMTxHash string `gorm:"column:update_evm_tx_hash;type:TEXT;not null"`
	Removed         bool   `gorm:"column:removed;default:false"`
}

func (*StorageProvider) TableName() string {
	return "storage_providers"
}

func (b *StorageProvider) ToSpEvent(e string) *SpEvent {
	return &SpEvent{
		SpID:      fmt.Sprintf("%d", b.SpID),
		Height:    b.CreateAt,
		TxHash:    b.CreateTxHash,
		EVMTxHash: b.CreateEVMTxHash,
		Event:     e,
	}
}

func NewSpEvent(spID string, height int64, txHash, evmTxHash, event string) *SpEvent {
	return &SpEvent{
		SpID:      spID,
		Height:    height,
		TxHash:    txHash,
		EVMTxHash: evmTxHash,
		Event:     event,
	}
}
