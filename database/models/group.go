package models

import (
	"time"

	"gorm.io/datatypes"
)

type Group struct {
	ID              uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	OwnerAddress    string         `gorm:"column:owner_address;type:TEXT;index:idx_owner"`
	GroupID         string         `gorm:"column:group_id;type:TEXT;index:idx_group_id;uniqueIndex:idx_account_group,priority:2"`
	GroupName       string         `gorm:"column:group_name;type:VARCHAR(63);index:idx_group_name"`
	SourceType      string         `gorm:"column:source_type;type:VARCHAR(63)"`
	Extra           string         `gorm:"column:extra;type:VARCHAR(512)"`
	AccountAddress  string         `gorm:"column:account_address;type:TEXT;uniqueIndex:idx_account_group,priority:1"`
	OperatorAddress string         `gorm:"column:operator_address;type:TEXT"`
	ExpirationTime  time.Time      `gorm:"column:expiration_time"`
	CreateAt        int64          `gorm:"column:create_at"`
	CreateTime      time.Time      `gorm:"column:create_time"`
	CreateTxHash    string         `gorm:"column:create_tx_hash;type:TEXT;not null"`
	CreateEVMTxHash string         `gorm:"column:create_evm_tx_hash;type:TEXT;not null"`
	UpdateAt        int64          `gorm:"column:update_at"`
	UpdateTime      time.Time      `gorm:"column:update_time"`
	UpdateTxHash    string         `gorm:"column:update_tx_hash;type:TEXT;not null"`
	UpdateEVMTxHash string         `gorm:"column:update_evm_tx_hash;type:TEXT;not null"`
	Removed         bool           `gorm:"column:removed;default:false"`
	Tags            datatypes.JSON `gorm:"column:tags;type:JSON"`
}

func (*Group) TableName() string {
	return "groups"
}

func (b *Group) ToBucketEvent(e string) *GroupEvent {
	return &GroupEvent{
		GroupID:   b.GroupID,
		Height:    b.CreateAt,
		TxHash:    b.CreateTxHash,
		EVMTxHash: b.CreateEVMTxHash,
		Event:     e,
	}
}

func NewGroupEvent(groupID string, height int64, txHash, evmTxHash, event string) *GroupEvent {
	return &GroupEvent{
		GroupID:   groupID,
		Height:    height,
		TxHash:    txHash,
		EVMTxHash: evmTxHash,
		Event:     event,
	}
}
