package models

type BucketEvent struct {
	ID        uint64 `gorm:"id;type:bigint(64);primaryKey;autoIncrement"`
	BucketID  string `gorm:"bucket_id;type:varchar(64);not null"`
	Height    int64  `gorm:"column:height;not null"`
	TxHash    string `gorm:"column:tx_hash;type:TEXT;not null"`
	EVMTxHash string `gorm:"column:evm_tx_hash;type:TEXT;not null"`
	Event     string `gorm:"column:event;type:TEXT;not null"`
}

func (*BucketEvent) TableName() string {
	return "bucket_events"
}

type ObjectEvent struct {
	ID        uint64 `gorm:"id;type:bigint(64);primaryKey;autoIncrement"`
	ObjectID  string `gorm:"object_id;type:varchar(64);not null"`
	Height    int64  `gorm:"column:height;not null"`
	TxHash    string `gorm:"column:tx_hash;type:TEXT;not null"`
	EVMTxHash string `gorm:"column:evm_tx_hash;type:TEXT;not null"`
	Event     string `gorm:"column:event;type:TEXT;not null"`
}

func (*ObjectEvent) TableName() string {
	return "object_events"
}

type GroupEvent struct {
	ID        uint64 `gorm:"id;type:bigint(64);primaryKey;autoIncrement"`
	GroupID   string `gorm:"group_id;type:varchar(64);not null"`
	Height    int64  `gorm:"column:height;not null"`
	TxHash    string `gorm:"column:tx_hash;type:TEXT;not null"`
	EVMTxHash string `gorm:"column:evm_tx_hash;type:TEXT;not null"`
	Event     string `gorm:"column:event;type:TEXT;not null"`
}

func (*GroupEvent) TableName() string {
	return "group_events"
}

type SpEvent struct {
	ID        uint64 `gorm:"id;type:bigint(64);primaryKey;autoIncrement"`
	SpID      string `gorm:"sp_id;type:varchar(64);not null"`
	Height    int64  `gorm:"column:height;not null"`
	TxHash    string `gorm:"column:tx_hash;type:TEXT;not null"`
	EVMTxHash string `gorm:"column:evm_tx_hash;type:TEXT;not null"`
	Event     string `gorm:"column:event;type:TEXT;not null"`
}

func (*SpEvent) TableName() string {
	return "sp_events"
}
