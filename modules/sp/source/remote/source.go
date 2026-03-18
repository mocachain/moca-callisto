package remote

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	sptypes "github.com/evmos/evmos/v12/x/sp/types"
	"github.com/forbole/bdjuno/v4/modules/sp/source"
	"github.com/forbole/juno/v5/node/remote"
)

var _ source.Source = &Source{}

// Source implements storage.Source using a remote node
type Source struct {
	*remote.Source
	Cli sptypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, cli sptypes.QueryClient) *Source {
	return &Source{
		Source: source,
		Cli:    cli,
	}
}

func (s Source) StorageProvider(height int64, id uint32) (sptypes.StorageProvider, error) {
	res, err := s.Cli.StorageProvider(
		remote.GetHeightRequestContext(s.Ctx, height),
		&sptypes.QueryStorageProviderRequest{Id: id},
	)
	if err != nil {
		return sptypes.StorageProvider{}, err
	}

	return *res.StorageProvider, nil
}

func (s Source) StorageProviders(height int64, pageRequest query.PageRequest) ([]*sptypes.StorageProvider, *query.PageResponse, error) {
	res, err := s.Cli.StorageProviders(
		remote.GetHeightRequestContext(s.Ctx, height),
		&sptypes.QueryStorageProvidersRequest{Pagination: &pageRequest},
	)
	if err != nil {
		return nil, nil, err
	}

	return res.Sps, res.Pagination, nil
}
