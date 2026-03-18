package source

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	sptypes "github.com/evmos/evmos/v12/x/sp/types"
)

type Source interface {
	StorageProvider(height int64, id uint32) (sptypes.StorageProvider, error)
	StorageProviders(height int64, pageRequest query.PageRequest) ([]*sptypes.StorageProvider, *query.PageResponse, error)
}
