package sp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "sp").Msg("parsing genesis")
	req := query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: false,
		Reverse:    false,
	}
	sps, pageResponse, err := m.source.StorageProviders(1, req)
	_ = pageResponse
	if err != nil {
		return err
	}
	ctx := context.Background()

	spModels := make([]*models.StorageProvider, 0, len(sps))
	for _, sp := range sps {
		s := &models.StorageProvider{
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
			CreateTxHash:    "",
			CreateAt:        1,
			UpdateAt:        1,
			UpdateTxHash:    "",
			Removed:         false,
		}
		spModels = append(spModels, s)
		k, v := m.db.CreateStorageProviderToSQL(ctx, s)
		if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
			return err
		}
	}

	for _, sp := range spModels {
		ek, ev := m.db.SaveSPEventToSQL(ctx, models.NewSpEvent(fmt.Sprintf("%d", sp.SpID), 1, "", "", EventCreateStorageProvider))
		if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: ek, Vars: ev}}); err != nil {
			return err
		}
	}
	return nil
}
