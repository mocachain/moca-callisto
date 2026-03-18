package virtualgroup

import (
	"context"
	"fmt"
	"time"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "virtualgroup").Msg("setting up periodic tasks")

	// Sync GVG families from chain every 5 minutes
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.SyncGVGFamiliesFromChain)
	}); err != nil {
		return fmt.Errorf("error while scheduling GVG families sync periodic operation: %s", err)
	}

	return nil
}

// SyncGVGFamiliesFromChain queries all GVG families from chain and syncs them to database
func (m *Module) SyncGVGFamiliesFromChain() error {
	ctx := context.Background()

	// Get current block height
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	// Query all GVG families from chain
	gvgFamilies, err := m.source.GlobalVirtualGroupFamilies(height, context.Background())
	if err != nil {
		log.Warn().Err(err).Str("module", "virtualgroup").Msg("failed to query global virtual group families from chain")
		return nil // Don't fail, just log and continue
	}

	if len(gvgFamilies) == 0 {
		return nil
	}

	log.Info().Int("count", len(gvgFamilies)).Str("module", "virtualgroup").Msg("found GVG families on chain, syncing to database")

	// Sync each GVG family to database
	for _, gvgf := range gvgFamilies {
		// Check if GVG family already exists
		existingGVGF, _ := m.db.GetGVGFamilyByID(ctx, gvgf.Id)

		vgfGroup := &models.GlobalVirtualGroupFamily{
			GlobalVirtualGroupFamilyID: gvgf.Id,
			PrimarySpID:                gvgf.PrimarySpId,
			VirtualPaymentAddress:      gvgf.VirtualPaymentAddress,
			GlobalVirtualGroupIDs:      models.ConvertUint32ToInt32Array(gvgf.GlobalVirtualGroupIds),
			UpdateAt:                   height,
			Removed:                    false,
		}

		// Preserve create info if GVG family already exists
		if existingGVGF != nil {
			vgfGroup.ID = existingGVGF.ID
			vgfGroup.CreateAt = existingGVGF.CreateAt
			vgfGroup.CreateTxHash = existingGVGF.CreateTxHash
			vgfGroup.CreateTime = existingGVGF.CreateTime
		} else {
			vgfGroup.CreateAt = height
			vgfGroup.CreateTxHash = ""
			vgfGroup.CreateEVMTxHash = ""
		}

		// Use SyncGVGFToSQL which only updates chain-queryable fields,
		// preserving event-sourced fields (update_tx_hash, update_evm_tx_hash, update_time)
		k, v := m.db.SyncGVGFToSQL(ctx, vgfGroup)
		if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
			log.Error().Err(err).Uint32("gvgf_id", gvgf.Id).Msg("failed to sync global virtual group family")
			continue
		}
		if existingGVGF == nil {
			log.Info().Uint32("gvgf_id", gvgf.Id).Msg("created global virtual group family from chain sync")
		}

		// Sync GVGs in this family
		gvgs, err := m.source.GlobalVirtualGroupByFamilyID(height, gvgf.Id)
		if err != nil {
			log.Warn().Err(err).Uint32("family_id", gvgf.Id).Msg("failed to get global virtual groups by family id")
			continue
		}

		for _, gvg := range gvgs {
			existingGVG, _ := m.db.GetGVGByID(ctx, gvg.Id)

			gvgGroup := &models.GlobalVirtualGroup{
				GlobalVirtualGroupID:  gvg.Id,
				FamilyID:              gvg.FamilyId,
				PrimarySpID:           gvg.PrimarySpId,
				SecondarySpIDs:        models.ConvertUint32ToInt32Array(gvg.SecondarySpIds),
				StoredSize:            gvg.StoredSize,
				VirtualPaymentAddress: gvg.VirtualPaymentAddress,
				TotalDeposit:          gvg.TotalDeposit.BigInt().String(),
				UpdateAt:              height,
				Removed:               false,
			}

			if existingGVG != nil {
				gvgGroup.ID = existingGVG.ID
				gvgGroup.CreateAt = existingGVG.CreateAt
				gvgGroup.CreateTxHash = existingGVG.CreateTxHash
				gvgGroup.CreateTime = existingGVG.CreateTime
			} else {
				gvgGroup.CreateAt = height
				gvgGroup.CreateTxHash = ""
				gvgGroup.CreateEVMTxHash = ""
				gvgGroup.CreateTime = time.Time{}
			}

			// Use SyncGVGToSQL which only updates chain-queryable fields,
			// preserving event-sourced fields (update_tx_hash, update_evm_tx_hash, update_time)
			k, v := m.db.SyncGVGToSQL(ctx, gvgGroup)
			if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
				log.Error().Err(err).Uint32("gvg_id", gvg.Id).Msg("failed to sync global virtual group")
				continue
			}
		}
	}

	return nil
}
