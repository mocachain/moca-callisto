package types

import "time"

type StorageGroupRow struct {
	GroupId    uint64 `db:"id"`
	GroupName  string `db:"group_name"`
	Owner      string `db:"owner"`
	SourceType string `db:"source_type"`
	Extra      string `db:"extra"`
	Height     int64  `db:"height"`
	Tags       string `db:"tags"`
}

func NewStorageGroupRow(
	groupId uint64,
	groupName string,
	owner string,
	sourceType string,
	extra string,
	height int64,
	tags string,
) StorageGroupRow {
	return StorageGroupRow{
		GroupId:    groupId,
		GroupName:  groupName,
		Owner:      owner,
		SourceType: sourceType,
		Extra:      extra,
		Height:     height,
		Tags:       tags,
	}
}
func (g StorageGroupRow) Equal(v StorageGroupRow) bool {
	return g.GroupId == v.GroupId &&
		g.GroupName == v.GroupName &&
		g.Owner == v.Owner &&
		g.SourceType == v.SourceType &&
		g.Extra == v.Extra &&
		g.Height == v.Height &&
		g.Tags == v.Tags
}

type BucketRow struct {
	BucketId                   uint64    `db:"id"`
	BucketName                 string    `db:"bucket_name"`
	Owner                      string    `db:"owner"`
	Visibility                 string    `db:"visibility"`
	SourceType                 string    `db:"source_type"`
	CreateAt                   time.Time `db:"create_at"`
	PaymentAddress             string    `db:"payment_address"`
	BucketStatus               string    `db:"bucket_status"`
	Tags                       string    `db:"tags"`
	ChargedReadQuota           uint64    `db:"charged_read_quota"`
	GlobalVirtualGroupFamilyId uint32    `db:"global_virtual_group_family_id"`
	SpAsDelegatedAgentDisabled bool      `db:"sp_as_delegated_agent_disabled"`
	Height                     int64     `db:"height"`
}

func NewBucketRow(
	bucketId uint64,
	bucketName string,
	owner string,
	visibility string,
	sourceType string,
	createAt time.Time,
	paymentAddress string,
	bucketStatus string,
	tags string,
	chargedReadQuota uint64,
	globalVirtualGroupFamilyId uint32,
	spAsDelegatedAgentDisabled bool,
	height int64,
) BucketRow {
	return BucketRow{
		BucketId:                   bucketId,
		BucketName:                 bucketName,
		Owner:                      owner,
		Visibility:                 visibility,
		SourceType:                 sourceType,
		CreateAt:                   createAt,
		PaymentAddress:             paymentAddress,
		BucketStatus:               bucketStatus,
		Tags:                       tags,
		ChargedReadQuota:           chargedReadQuota,
		GlobalVirtualGroupFamilyId: globalVirtualGroupFamilyId,
		SpAsDelegatedAgentDisabled: spAsDelegatedAgentDisabled,
		Height:                     height,
	}
}

func (b BucketRow) Equal(v BucketRow) bool {
	return b.BucketId == v.BucketId &&
		b.BucketName == v.BucketName &&
		b.Owner == v.Owner &&
		b.Visibility == v.Visibility &&
		b.SourceType == v.SourceType &&
		b.CreateAt.Equal(v.CreateAt) &&
		b.PaymentAddress == v.PaymentAddress &&
		b.BucketStatus == v.BucketStatus &&
		b.Tags == v.Tags &&
		b.ChargedReadQuota == v.ChargedReadQuota &&
		b.GlobalVirtualGroupFamilyId == v.GlobalVirtualGroupFamilyId &&
		b.SpAsDelegatedAgentDisabled == v.SpAsDelegatedAgentDisabled &&
		b.Height == v.Height
}

type ObjectRow struct {
	ObjectId            uint64    `db:"id"`
	ObjectName          string    `db:"object_name"`
	BucketName          string    `db:"bucket_name"`
	Owner               string    `db:"owner"`
	Creator             string    `db:"creator"`
	PayloadSize         uint64    `db:"payload_size"`
	Visibility          string    `db:"visibility"`
	ContentType         string    `db:"content_type"`
	ObjectStatus        string    `db:"object_status"`
	RedundancyType      string    `db:"redundancy_type"`
	SourceType          string    `db:"source_type"`
	Checksums           string    `db:"checksums"`
	Tags                string    `db:"tags"`
	IsUpdating          bool      `db:"is_updating"`
	CreateAt            time.Time `db:"create_at"`
	UpdatedAt           time.Time `db:"update_at"`
	UpdatedBy           string    `db:"update_by"`
	Version             int64     `db:"version"`
	LocalVirtualGroupId uint32    `db:"local_virtual_group_id"`
	Height              int64     `db:"height"`
}

func NewObjectsRow(
	objectId uint64,
	objectName string,
	bucketName string,
	owner string,
	creator string,
	payloadSize uint64,
	visibility string,
	contentType string,
	objectStatus string,
	redundancyType string,
	sourceType string,
	checksums string,
	tags string,
	isUpdating bool,
	createAt time.Time,
	updateAt time.Time,
	updatedBy string,
	version int64,
	localVirtualGroupId uint32,
	height int64,
) ObjectRow {
	return ObjectRow{
		ObjectId:            objectId,
		ObjectName:          objectName,
		BucketName:          bucketName,
		Owner:               owner,
		Creator:             creator,
		PayloadSize:         payloadSize,
		Visibility:          visibility,
		ContentType:         contentType,
		ObjectStatus:        objectStatus,
		RedundancyType:      redundancyType,
		SourceType:          sourceType,
		Checksums:           checksums,
		Tags:                tags,
		IsUpdating:          isUpdating,
		CreateAt:            createAt,
		UpdatedAt:           updateAt,
		UpdatedBy:           updatedBy,
		Version:             version,
		LocalVirtualGroupId: localVirtualGroupId,
		Height:              height,
	}
}

func (o ObjectRow) Equal(v ObjectRow) bool {
	return o.ObjectId == v.ObjectId &&
		o.ObjectName == v.ObjectName &&
		o.BucketName == v.BucketName &&
		o.Owner == v.Owner &&
		o.Creator == v.Creator &&
		o.PayloadSize == v.PayloadSize &&
		o.Visibility == v.Visibility &&
		o.ContentType == v.ContentType &&
		o.ObjectStatus == v.ObjectStatus &&
		o.RedundancyType == v.RedundancyType &&
		o.SourceType == v.SourceType &&
		o.Checksums == v.Checksums &&
		o.Tags == v.Tags &&
		o.IsUpdating == v.IsUpdating &&
		o.CreateAt.Equal(v.CreateAt) &&
		o.UpdatedAt.Equal(v.UpdatedAt) &&
		o.UpdatedBy == v.UpdatedBy &&
		o.Version == v.Version &&
		o.LocalVirtualGroupId == v.LocalVirtualGroupId &&
		o.Height == v.Height
}
