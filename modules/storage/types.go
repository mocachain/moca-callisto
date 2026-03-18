package storage

type createGroupEvent struct {
	GroupId    string `json:"group_id"`
	GroupName  string `json:"group_name"`
	Owner      string `json:"owner"`
	SourceType string `json:"source_type"`
	Extra      string `json:"extra"`
}

type deleteGroupEvent struct {
	GroupId   string `json:"group_id"`
	GroupName string `json:"group_name"`
	Owner     string `json:"owner"`
}

type createBucketEvent struct {
	BucketId                   string `json:"bucket_id"`
	BucketName                 string `json:"bucket_name"`
	Owner                      string `json:"owner"`
	Visibility                 string `json:"visibility"`
	CreateAt                   int64  `json:"create_at"`
	SourceType                 string `json:"source_type"`
	ChargedReadQuota           uint64 `json:"charged_read_quota"`
	PaymentAddress             string `json:"payment_address"`
	PrimarySpId                uint32 `json:"primary_sp_id"`
	GlobalVirtualGroupFamilyId uint32 `json:"global_virtual_group_family_id"`
	Status                     string `json:"status"`
}

type deleteBucketEvent struct {
	BucketId                   string `json:"bucket_id"`
	BucketName                 string `json:"bucket_name"`
	Owner                      string `json:"owner"`
	Operator                   string `json:"operator"`
	GlobalVirtualGroupFamilyId uint32 `json:"global_virtual_group_family_id"`
}

type createObjectEvent struct {
	ObjectId            string `json:"object_id"`
	ObjectName          string `json:"object_name"`
	BucketId            string `json:"bucket_id"`
	BucketName          string `json:"bucket_name"`
	Owner               string `json:"owner"`
	Creator             string `json:"creator"`
	PrimarySpId         uint32 `json:"primary_sp_id"`
	PayloadSize         uint64 `json:"payload_size"`
	Visibility          string `json:"visibility"`
	ContentType         string `json:"content_type"`
	CreateAt            int64  `json:"create_at"`
	Status              string `json:"status"`
	RedundancyType      string `json:"redundancy_type"`
	SourceType          string `json:"source_type"`
	Checksum            string `json:"checksum"`
	LocalVirtualGroupId uint32 `json:"local_virtual_group_id"`
}

type deleteObjectEvent struct {
	ObjectId    string `json:"object_id"`
	ObjectName  string `json:"object_name"`
	BucketName  string `json:"bucket_name"`
	Operator    string `json:"operator"`
	PrimarySpId uint32 `json:"primary_sp_id"`
}
