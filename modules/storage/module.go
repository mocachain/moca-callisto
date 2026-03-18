package storage

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/storage/source"
	vgsource "github.com/forbole/bdjuno/v4/modules/virtualgroup/source"
	"github.com/forbole/juno/v5/modules"
)

var (
	_ modules.Module      = &Module{}
	_ modules.BlockModule = &Module{}
)

// Module represents the x/distr module
type Module struct {
	source   source.Source
	vgSource vgsource.Source
	cdc      codec.Codec
	db       *database.DB
}

// NewModule returns a new Module instance
func NewModule(source source.Source, vgSource vgsource.Source, cdc codec.Codec, db *database.DB) *Module {
	return &Module{
		source:   source,
		vgSource: vgSource,
		cdc:      cdc,
		db:       db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "storage"
}
