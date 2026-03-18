package types

import "time"

type StorageGroup struct {
	GroupId    uint64
	GroupName  string
	Owner      string
	SourceType string
	Extra      string
	Height     int64
	Tags       string
}

func NewStorageGroup(
	groupId uint64,
	groupName string,
	owner string,
	sourceType string,
	extra string,
	height int64,
	tags string,
) StorageGroup {
	return StorageGroup{
		GroupId:    groupId,
		GroupName:  groupName,
		Owner:      owner,
		SourceType: sourceType,
		Extra:      extra,
		Height:     height,
		Tags:       tags,
	}
}

func (g StorageGroup) Equal(v StorageGroup) bool {
	return g.GroupId == v.GroupId &&
		g.GroupName == v.GroupName &&
		g.Owner == v.Owner &&
		g.SourceType == v.SourceType &&
		g.Extra == v.Extra &&
		g.Height == v.Height &&
		g.Tags == v.Tags
}

type Bucket struct {
	BucketId                   uint64
	BucketName                 string
	Owner                      string
	Visibility                 string
	SourceType                 string
	CreateAt                   time.Time
	PaymentAddress             string
	BucketStatus               string
	Tags                       string
	ChargedReadQuota           uint64
	GlobalVirtualGroupFamilyId uint32
	SpAsDelegatedAgentDisabled bool
	Height                     int64
}

func NewBucket(
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
) Bucket {
	return Bucket{
		BucketId:                   bucketId,
		BucketName:                 bucketName,
		Owner:                      owner,
		Visibility:                 visibility,
		SourceType:                 sourceType,
		PaymentAddress:             paymentAddress,
		BucketStatus:               bucketStatus,
		Tags:                       tags,
		ChargedReadQuota:           chargedReadQuota,
		GlobalVirtualGroupFamilyId: globalVirtualGroupFamilyId,
		SpAsDelegatedAgentDisabled: spAsDelegatedAgentDisabled,
		Height:                     height,
	}
}

func (b Bucket) Equal(v Bucket) bool {
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

type Object struct {
	ObjectId            uint64
	ObjectName          string
	BucketName          string
	Owner               string
	Creator             string
	PayloadSize         uint64
	Visibility          string
	ContentType         string
	ObjectStatus        string
	RedundancyType      string
	SourceType          string
	Checksums           string
	Tags                string
	IsUpdating          bool
	CreateAt            time.Time
	UpdatedAt           time.Time
	UpdatedBy           string
	Version             int64
	LocalVirtualGroupId uint32
	Height              int64
}

func NewObject(
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
) Object {
	return Object{
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

func (o Object) Equal(v Object) bool {
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
