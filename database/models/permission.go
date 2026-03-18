package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/lib/pq"
)

type Permission struct {
	ID             uint64    `gorm:"column:id;type:bigint(64);primaryKey;autoIncrement"`
	PrincipalType  int32     `gorm:"column:principal_type;type:int;uniqueIndex:idx_policy,priority:1"`
	PrincipalValue string    `gorm:"column:principal_value;type:varchar(128);uniqueIndex:idx_policy,priority:2"`
	ResourceType   string    `gorm:"column:resource_type;type:varchar(64);uniqueIndex:idx_policy,priority:3"`
	ResourceID     string    `gorm:"column:resource_id;type:TEXT;uniqueIndex:idx_policy,priority:4"`
	PolicyID       string    `gorm:"column:policy_id;type:TEXT;index:idx_policy_id"`
	CreateTime     time.Time `gorm:"column:create_time;type:timestamp"`
	UpdateTime     time.Time `gorm:"column:update_time;type:timestamp"`
	ExpirationTime time.Time `gorm:"column:expiration_time;type:timestamp"`
	Removed        bool      `gorm:"column:removed"`
}

func (p Permission) TableName() string {
	return "permission"
}

type Statements struct {
	ID             uint64         `gorm:"id;type:bigint(64);primaryKey;autoIncrement"`
	PolicyID       common.Hash    `gorm:"policy_id;type:BINARY(32);index:idx_policy_id"`
	Effect         string         `gorm:"effect;type:varchar(32)"`
	ActionValue    int            `gorm:"action_value;type:int"`
	Resources      pq.StringArray `gorm:"resources;type:text"`
	ExpirationTime time.Time      `gorm:"expiration_time;type:bigint(64)"`
	LimitSize      uint64         `gorm:"limit_size;type:bigint(64)"`
	Removed        bool           `gorm:"removed;"`
}

func (s Statements) TableName() string {
	return "statements"
}

type GroupMember struct {
	ID      uint64         `gorm:"id;type:bigint(64);primaryKey;autoIncrement"`
	GroupID uint64         `gorm:"id;type:bigint(64);index:idx_group_id"`
	Account common.Address `gorm:"account;type:BINARY(20);index:idx_account"`
	Removed bool           `gorm:"removed"`
}

func (g GroupMember) TableName() string {
	return "permission"
}
