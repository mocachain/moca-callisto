package remote

import (
	permissiontypes "github.com/evmos/evmos/v12/x/permission/types"
	"github.com/forbole/bdjuno/v4/modules/storage/source"
	"github.com/evmos/evmos/v12/x/storage/types"
	vgtypes "github.com/evmos/evmos/v12/x/virtualgroup/types"
	"github.com/forbole/juno/v5/node/remote"
)

var _ source.Source = &Source{}

// Source implements storage.Source using a remote node
type Source struct {
	*remote.Source
	Cli types.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, cli types.QueryClient) *Source {
	return &Source{
		Source: source,
		Cli:    cli,
	}
}

// HeadBucket implements storage source.Source
func (s Source) HeadBucket(height int64, bucketName string) (types.BucketInfo, types.BucketExtraInfo, error) {
	res, err := s.Cli.HeadBucket(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadBucketRequest{BucketName: bucketName},
	)
	if err != nil {
		return types.BucketInfo{}, types.BucketExtraInfo{}, err
	}

	return *res.BucketInfo, *res.ExtraInfo, nil
}

func (s Source) HeadBucketById(height int64, bucketId string) (types.BucketInfo, types.BucketExtraInfo, error) {
	res, err := s.Cli.HeadBucketById(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadBucketByIdRequest{BucketId: bucketId},
	)
	if err != nil {
		return types.BucketInfo{}, types.BucketExtraInfo{}, err
	}

	return *res.BucketInfo, *res.ExtraInfo, nil
}

func (s Source) HeadBucketExtra(height int64, bucketName string) (types.InternalBucketInfo, error) {
	res, err := s.Cli.HeadBucketExtra(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadBucketExtraRequest{BucketName: bucketName},
	)
	if err != nil {
		return types.InternalBucketInfo{}, err
	}

	return *res.ExtraInfo, nil
}

func (s Source) HeadGroup(height int64, groupOwner, groupName string) (types.GroupInfo, error) {
	res, err := s.Cli.HeadGroup(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadGroupRequest{
			GroupOwner: groupOwner,
			GroupName:  groupName,
		},
	)
	if err != nil {
		return types.GroupInfo{}, err
	}

	return *res.GroupInfo, nil
}

func (s Source) HeadGroupMember(height int64, member, groupOwner, groupName string) (permissiontypes.GroupMember, error) {
	res, err := s.Cli.HeadGroupMember(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadGroupMemberRequest{
			Member:     member,
			GroupOwner: groupOwner,
			GroupName:  groupName,
		},
	)
	if err != nil {
		return permissiontypes.GroupMember{}, err
	}

	return *res.GroupMember, nil
}

func (s Source) HeadObject(height int64, bucketName, objectName string) (types.ObjectInfo, vgtypes.GlobalVirtualGroup, error) {
	res, err := s.Cli.HeadObject(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadObjectRequest{BucketName: bucketName, ObjectName: objectName},
	)
	if err != nil {
		return types.ObjectInfo{}, vgtypes.GlobalVirtualGroup{}, err
	}

	return *res.ObjectInfo, *res.GlobalVirtualGroup, nil
}

func (s Source) HeadObjectById(height int64, objectId string) (types.ObjectInfo, vgtypes.GlobalVirtualGroup, error) {
	res, err := s.Cli.HeadObjectById(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadObjectByIdRequest{ObjectId: objectId},
	)
	if err != nil {
		return types.ObjectInfo{}, vgtypes.GlobalVirtualGroup{}, err
	}

	return *res.ObjectInfo, *res.GlobalVirtualGroup, nil
}

func (s Source) HeadShadowObject(height int64, bucketName, objectName string) (types.ShadowObjectInfo, error) {
	res, err := s.Cli.HeadShadowObject(
		remote.GetHeightRequestContext(s.Ctx, height),
		&types.QueryHeadShadowObjectRequest{BucketName: bucketName, ObjectName: objectName},
	)
	if err != nil {
		return types.ShadowObjectInfo{}, err
	}

	return *res.ObjectInfo, nil
}
