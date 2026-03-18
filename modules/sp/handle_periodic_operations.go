package sp

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "sp").Msg("setting up periodic tasks")

	// Sync storage providers from chain every 5 minutes
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.SyncStorageProvidersFromChain)
	}); err != nil {
		return fmt.Errorf("error while scheduling SP sync periodic operation: %s", err)
	}

	return nil
}

// SyncStorageProvidersFromChain queries all storage providers from chain and syncs them to database
func (m *Module) SyncStorageProvidersFromChain() error {
	ctx := context.Background()

	// Get current block height
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	log.Debug().Str("module", "sp").Int64("height", height).
		Msg("syncing storage providers from chain")

	// Query all storage providers from chain
	req := query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: false,
		Reverse:    false,
	}

	sps, _, err := m.source.StorageProviders(height, req)
	if err != nil {
		log.Warn().Err(err).Str("module", "sp").Msg("failed to query storage providers from chain")
		return nil // Don't fail, just log and continue
	}

	if len(sps) == 0 {
		log.Debug().Str("module", "sp").Msg("no storage providers found on chain")
		return nil
	}

	log.Info().Int("count", len(sps)).Str("module", "sp").Msg("found storage providers on chain, syncing to database")

	// Sync each SP to database
	for _, sp := range sps {
		// Check if SP already exists to preserve create info
		existingSp, _ := m.db.GetStorageProviderBySpID(ctx, sp.Id)

		spModel := &models.StorageProvider{
			SpID:            sp.Id,
			OperatorAddress: sp.OperatorAddress,
			FundingAddress:  sp.FundingAddress,
			SealAddress:     sp.SealAddress,
			ApprovalAddress: sp.ApprovalAddress,
			GcAddress:       sp.GcAddress,
			TotalDeposit:    sp.TotalDeposit.BigInt().String(),
			Status:          sp.Status.String(),
			Endpoint:        sp.Endpoint,
			Moniker:         sp.Description.Moniker,
			Identity:        sp.Description.Identity,
			Website:         sp.Description.Website,
			SecurityContact: sp.Description.SecurityContact,
			Details:         sp.Description.Details,
			BlsKey:          hex.EncodeToString(sp.BlsKey),
			UpdateAt:        height,
			Removed:         false,
		}

		// Preserve create info if SP already exists
		if existingSp != nil {
			spModel.ID = existingSp.ID
			spModel.CreateAt = existingSp.CreateAt
			spModel.CreateTxHash = existingSp.CreateTxHash
			spModel.CreateEVMTxHash = existingSp.CreateEVMTxHash
		} else {
			spModel.CreateAt = height
			spModel.CreateTxHash = "" // We don't have tx hash from query
			spModel.CreateEVMTxHash = ""
		}

		// Use SyncStorageProviderToSQL which only updates chain-queryable fields,
		// preserving event-sourced fields (prices, quotas, timestamps)
		k, v := m.db.SyncStorageProviderToSQL(ctx, spModel)
		if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
			log.Error().Err(err).Uint32("sp_id", sp.Id).Msg("failed to sync storage provider")
			continue
		}
		if existingSp == nil {
			log.Info().Uint32("sp_id", sp.Id).Msg("created storage provider from chain sync")
		} else {
			log.Debug().Uint32("sp_id", sp.Id).Msg("updated storage provider from chain sync")
		}
	}

	return nil
}
