package models

import (
	"time"

	"github.com/shopspring/decimal"
	datatypes "gorm.io/datatypes"
)

type Bucket struct {
	ID                         uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	BucketID                   string          `gorm:"column:bucket_id;type:TEXT;uniqueIndex:idx_bucket_id"`
	BucketName                 string          `gorm:"column:bucket_name;type:VARCHAR(64)"`
	OwnerAddress               string          `gorm:"column:owner_address;type:TEXT;index:idx_owner"`
	PaymentAddress             string          `gorm:"column:payment_address;type:TEXT"`
	GlobalVirtualGroupFamilyId uint32          `gorm:"column:global_virtual_group_family_id;index:idx_vgf_id"`
	OperatorAddress            string          `gorm:"column:operator_address;type:TEXT"`
	SourceType                 string          `gorm:"column:source_type;type:VARCHAR(50)"`
	ChargedReadQuota           uint64          `gorm:"column:charged_read_quota"`
	Visibility                 string          `gorm:"column:visibility;type:VARCHAR(50)"`
	Status                     string          `gorm:"column:status;type:VARCHAR(64)"`
	DeleteAt                   int64           `gorm:"column:delete_at"`
	DeleteReason               string          `gorm:"column:delete_reason;type:VARCHAR(256)"`
	StorageSize                decimal.Decimal `gorm:"column:storage_size;type:DECIMAL(65, 0);not null"`
	ChargeSize                 decimal.Decimal `gorm:"column:charge_size;type:DECIMAL(65, 0);not null"`
	CreateAt                   int64           `gorm:"column:create_at"`
	CreateTime                 time.Time       `gorm:"column:create_time"`
	CreateTxHash               string          `gorm:"column:create_tx_hash;type:TEXT;not null"`
	CreateEVMTxHash            string          `gorm:"column:create_evm_tx_hash;type:TEXT;not null"`
	UpdateAt                   int64           `gorm:"column:update_at"`
	UpdateTxHash               string          `gorm:"column:update_tx_hash;type:TEXT;not null"`
	UpdateEVMTxHash            string          `gorm:"column:update_evm_tx_hash;type:TEXT;not null"`
	UpdateTime                 time.Time       `gorm:"column:update_time"`
	Removed                    bool            `gorm:"column:removed;default:false"`
	OffChainStatus             int             `gorm:"column:off_chain_status;type:INT;not null;default:0"`
	SpAsDelegatedAgentDisabled bool            `gorm:"column:sp_as_delegated_agent_disabled;default:false"`
	Tags                       datatypes.JSON  `gorm:"column:tags;type:JSON"`
}

func (*Bucket) TableName() string {
	return "buckets"
}

func (b *Bucket) ToBucketEvent(e string) *BucketEvent {
	return &BucketEvent{
		BucketID:  b.BucketID,
		Height:    b.CreateAt,
		TxHash:    b.CreateTxHash,
		EVMTxHash: b.CreateEVMTxHash,
		Event:     e,
	}
}

func NewBucketEvent(bucketID string, height int64, txHash, evmTxHash, event string) *BucketEvent {
	return &BucketEvent{
		BucketID:  bucketID,
		Height:    height,
		TxHash:    txHash,
		EVMTxHash: evmTxHash,
		Event:     event,
	}
}
