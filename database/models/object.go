package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Object struct {
	ID                  uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	BucketID            string         `gorm:"column:bucket_id;type:TEXT;index:idx_bucket_id"`
	BucketName          string         `gorm:"column:bucket_name;type:VARCHAR(64);index:idx_bucket_name_object_name,priority:1"`
	ObjectID            string         `gorm:"column:object_id;type:TEXT;uniqueIndex:idx_object_id"`
	ObjectName          string         `gorm:"column:object_name;type:VARCHAR(1024);index:idx_bucket_name_object_name,length:512,priority:2"`
	CreatorAddress      string         `gorm:"column:creator_address;type:TEXT"`
	OwnerAddress        string         `gorm:"column:owner_address;type:TEXT;index:idx_owner"`
	LocalVirtualGroupID uint32         `gorm:"column:local_virtual_group_id;index:idx_lvg_id"`
	OperatorAddress     string         `gorm:"column:operator_address;type:TEXT"`
	PayloadSize         uint64         `gorm:"column:payload_size"`
	Visibility          string         `gorm:"column:visibility;type:VARCHAR(50)"`
	ContentType         string         `gorm:"column:content_type"`
	Status              string         `gorm:"column:status;type:VARCHAR(50)"`
	RedundancyType      string         `gorm:"column:redundancy_type;type:VARCHAR(50)"`
	SourceType          string         `gorm:"column:source_type;type:VARCHAR(50)"`
	CheckSums           pq.StringArray `gorm:"column:checksums;type:TEXT[]"`
	DeleteAt            int64          `gorm:"column:delete_at"`
	DeleteReason        string         `gorm:"column:delete_reason;type:VARCHAR(256)"`
	CreateAt            int64          `gorm:"column:create_at"`
	CreateTxHash        string         `gorm:"column:create_tx_hash;type:TEXT;not null"`
	CreateEVMTxHash     string         `gorm:"column:create_evm_tx_hash;type:TEXT;not null"`
	CreateTime          time.Time      `gorm:"column:create_time"`
	UpdateAt            int64          `gorm:"column:update_at;index:idx_update_at"`
	UpdateTxHash        string         `gorm:"column:update_tx_hash;type:TEXT;not null"`
	UpdateEVMTxHash     string         `gorm:"column:update_evm_tx_hash;type:TEXT;not null"`
	SealedTxHash        string         `gorm:"column:sealed_tx_hash;type:TEXT"`
	SealedEvmTxHash     string         `gorm:"column:sealed_evm_tx_hash;type:TEXT"`
	UpdateTime          time.Time      `gorm:"column:update_time"`
	Removed             bool           `gorm:"column:removed;default:false"`
	Tags                datatypes.JSON `gorm:"column:tags;type:JSON"`
	IsUpdating          bool           `gorm:"column:is_updating"`
	ContentUpdatedTime  time.Time      `gorm:"column:content_updated_time"`
	Updater             string         `gorm:"column:updater;type:TEXT"`
	Version             int64          `gorm:"column:version"`
}

func (*Object) TableName() string {
	return "objects"
}

func (b *Object) ToObjectEvent(e string) *ObjectEvent {
	return &ObjectEvent{
		ObjectID:  b.ObjectID,
		Height:    b.CreateAt,
		TxHash:    b.CreateTxHash,
		EVMTxHash: b.CreateEVMTxHash,
		Event:     e,
	}
}

func NewObjectEvent(objectID string, height int64, txHash, evmTxHash, event string) *ObjectEvent {
	return &ObjectEvent{
		ObjectID:  objectID,
		Height:    height,
		TxHash:    txHash,
		EVMTxHash: evmTxHash,
		Event:     event,
	}
}
