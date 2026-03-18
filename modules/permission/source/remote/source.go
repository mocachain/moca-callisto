package remote

import (
	"github.com/forbole/bdjuno/v4/modules/permission/source"
	"github.com/evmos/evmos/v12/x/permission/types"
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
