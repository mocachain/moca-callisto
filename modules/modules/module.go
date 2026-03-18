package modules

import (
	"github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/types/config"

	"github.com/forbole/bdjuno/v4/database"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

type Module struct {
	cfg config.ChainConfig
	db  *database.DB
}

// NewModule returns a new Module instance
func NewModule(cfg config.ChainConfig, db *database.DB) *Module {
	return &Module{
		cfg: cfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "modules"
}
