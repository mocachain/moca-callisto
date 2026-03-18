package source

import (
	permissiontypes "github.com/evmos/evmos/v12/x/permission/types"
	"github.com/evmos/evmos/v12/x/storage/types"
	vgtypes "github.com/evmos/evmos/v12/x/virtualgroup/types"
)

type Source interface {
	HeadBucket(height int64, bucketName string) (types.BucketInfo, types.BucketExtraInfo, error)
	HeadBucketById(height int64, bucketId string) (types.BucketInfo, types.BucketExtraInfo, error)
	HeadBucketExtra(height int64, bucketName string) (types.InternalBucketInfo, error)

	HeadGroup(height int64, groupOwner, groupName string) (types.GroupInfo, error)
	HeadGroupMember(height int64, member, groupOwner, groupName string) (permissiontypes.GroupMember, error)

	HeadObject(height int64, bucketName, objectName string) (types.ObjectInfo, vgtypes.GlobalVirtualGroup, error)
	HeadObjectById(height int64, objectId string) (types.ObjectInfo, vgtypes.GlobalVirtualGroup, error)
	HeadShadowObject(height int64, bucketName, objectName string) (types.ShadowObjectInfo, error)
}
