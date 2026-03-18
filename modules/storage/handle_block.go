package storage

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/database"
	junotypes "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
)

const (
	EventTypeEthereumTx = "ethereum_tx"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, txs []*junotypes.Tx, vals *tmctypes.ResultValidators,
) error {
	ctx := context.Background()
	statements, err := m.ExportEventsInTxs(ctx, block, txs)
	if err != nil {
		return err
	}
	return m.db.ExecuteStatements(statements)
}

// ExportEventsInTxs accepts a slice of events in tx in order to save in database.
func (m Module) ExportEventsInTxs(ctx context.Context, block *tmctypes.ResultBlock, txs []*junotypes.Tx) ([]database.SQLStatement, error) {
	allSQL := make([]database.SQLStatement, 0)
	for _, tx := range txs {
		sqls, err := m.ExtractEvent(ctx, block, tx)
		if err != nil {
			log.Err(err)
			continue
		}
		allSQL = append(allSQL, sqls...)
	}
	return allSQL, nil
}

// ExtractEvent accepts the transaction and handles events contained inside the transaction.
func (m *Module) ExtractEvent(ctx context.Context, block *tmctypes.ResultBlock, tx *junotypes.Tx) ([]database.SQLStatement, error) {
	txHash := tx.TxHash
	evmTxHash := findEVMTxHash(tx.Events)
	allSQL := make([]database.SQLStatement, 0)
	for _, event := range tx.Events {
		e := sdk.Event(event)
		h := m.getExtractEventFunc(e)
		if h == nil {
			continue
		}
		raw, err := h(ctx, block, txHash, evmTxHash, e)
		if err != nil {
			log.Err(err)
			continue
		}
		switch v := raw.(type) {
		case nil:
			continue
		case map[string][]interface{}:
			for k, args := range v {
				allSQL = append(allSQL, database.SQLStatement{SQL: k, Vars: args})
			}
		case []database.SQLStatement:
			allSQL = append(allSQL, v...)
		default:
			log.Warn().Str("event_type", e.Type).Msgf("unsupported SQL payload type: %T", raw)
		}
	}
	return allSQL, nil
}

type ExtractFunc func(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, event sdk.Event) (interface{}, error)

func (m *Module) getExtractEventFunc(event sdk.Event) ExtractFunc {
	switch {
	case BucketEvents[event.Type]:
		return m.ExtractBucketEventStatements
	case ObjectEvents[event.Type]:
		return m.ExtractObjectEventStatements
	case GroupEvents[event.Type]:
		return m.ExtractGroupEventStatements
	case TagEvents[event.Type]:
		return m.ExtractTagEventStatements
	default:
		return nil
	}
}

func findEVMTxHash(events []abci.Event) string {
	for _, event := range events {
		if event.Type != EventTypeEthereumTx {
			continue
		}
		for _, att := range event.Attributes {
			if att.Key == "ethereumTxHash" {
				return att.Value
			}
		}
	}
	return ""
}
