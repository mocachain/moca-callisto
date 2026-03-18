package virtualgroup

import (
	"context"
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "virtualgroup").Msg("parsing genesis")
	
	// Query GVG families from chain at genesis height
	gvgFamilies, err := m.source.GlobalVirtualGroupFamilies(1, context.Background())
	if err != nil {
		log.Warn().Err(err).Str("module", "virtualgroup").Msg("failed to query global virtual group families at genesis height, trying current height")
		currentGVGFamilies, currentErr := m.source.GlobalVirtualGroupFamilies(0, context.Background())
		if currentErr != nil {
			log.Warn().Err(currentErr).Str("module", "virtualgroup").Msg("failed to query global virtual group families at current height, will try to sync from events")
			return nil
		}
		gvgFamilies = currentGVGFamilies
		log.Info().Int("count", len(gvgFamilies)).Str("module", "virtualgroup").Msg("queried global virtual group families from current height")
	}
	
	if len(gvgFamilies) == 0 {
		log.Warn().Str("module", "virtualgroup").Msg("no global virtual group families found at query time, will try to sync from events")
		return nil
	}
	
	log.Info().Int("count", len(gvgFamilies)).Str("module", "virtualgroup").Msg("found global virtual group families in genesis")
	ctx := context.Background()
	
	for _, gvgf := range gvgFamilies {
		vgfGroup := &models.GlobalVirtualGroupFamily{
			GlobalVirtualGroupFamilyID: gvgf.Id,
			PrimarySpID:                gvgf.PrimarySpId,
			VirtualPaymentAddress:      gvgf.VirtualPaymentAddress,
			GlobalVirtualGroupIDs:      models.ConvertUint32ToInt32Array(gvgf.GlobalVirtualGroupIds),
			CreateAt:                   1,
			CreateTxHash:               "",
			CreateTime:                 doc.GenesisTime,
			UpdateAt:                   1,
			UpdateTxHash:               "",
			UpdateTime:                 doc.GenesisTime,
			Removed:                    false,
		}
		k, v := m.db.SaveGVGFToSQL(ctx, vgfGroup)
		if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
			log.Error().Err(err).Uint32("gvgf_id", gvgf.Id).Msg("failed to save global virtual group family from genesis")
			continue
		}
		
		// Also sync GVGs in this family
		gvgs, err := m.source.GlobalVirtualGroupByFamilyID(1, gvgf.Id)
		if err != nil {
			log.Warn().Err(err).Uint32("gvgf_id", gvgf.Id).Msg("failed to get global virtual groups by family id at genesis")
			// Try current height
			gvgs, err = m.source.GlobalVirtualGroupByFamilyID(0, gvgf.Id)
			if err != nil {
				log.Warn().Err(err).Uint32("gvgf_id", gvgf.Id).Msg("failed to get global virtual groups by family id at current height")
				continue
			}
		}
		
		for _, gvg := range gvgs {
			gvgGroup := &models.GlobalVirtualGroup{
				GlobalVirtualGroupID:  gvg.Id,
				FamilyID:              gvg.FamilyId,
				PrimarySpID:           gvg.PrimarySpId,
				SecondarySpIDs:         models.ConvertUint32ToInt32Array(gvg.SecondarySpIds),
				StoredSize:            gvg.StoredSize,
				VirtualPaymentAddress: gvg.VirtualPaymentAddress,
				TotalDeposit:          gvg.TotalDeposit.BigInt().String(),
				CreateAt:              1,
				CreateTxHash:          "",
				CreateTime:            doc.GenesisTime,
				UpdateAt:              1,
				UpdateTxHash:          "",
				UpdateTime:            doc.GenesisTime,
				Removed:               false,
			}
			k, v := m.db.SaveGVGToSQL(ctx, gvgGroup)
			if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
				log.Error().Err(err).Uint32("gvg_id", gvg.Id).Msg("failed to save global virtual group from genesis")
				continue
			}
		}
	}
	
	return nil
}

