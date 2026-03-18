package storage

import (
	"context"
	"encoding/json"
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/forbole/bdjuno/v4/database/models"
	storagetypes "github.com/evmos/evmos/v12/x/storage/types"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/forbole/bdjuno/v4/types/resource"
	"github.com/rs/zerolog/log"
)

var EventSetTag = proto.MessageName(&storagetypes.EventSetTag{})

var TagEvents = map[string]bool{
	EventSetTag: true,
}

func (m *Module) ExtractTagEventStatements(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, event sdk.Event) (interface{}, error) {
	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case EventSetTag:
		setTag, ok := typedEvent.(*storagetypes.EventSetTag)
		if !ok {
			return nil, errors.New("create object event assert error")
		}
		return m.handleSetTags(ctx, block, txHash, evmTxHash, setTag), nil
	default:
		return nil, errors.New("event type not found")
	}
}

func (m *Module) handleSetTags(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, e *storagetypes.EventSetTag) map[string][]interface{} {
	var grn types.GRN
	if err := grn.ParseFromString(e.Resource, true); err != nil {
		return nil
	}
	var (
		k string
		v []interface{}
	)
	switch grn.ResourceType() {
	case resource.RESOURCE_TYPE_BUCKET:
		b := &models.Bucket{
			BucketID:        e.Id.BigInt().String(),
			UpdateAt:        block.Block.Height,
			UpdateTxHash:    txHash,
			UpdateEVMTxHash: evmTxHash,
			UpdateTime:      block.Block.Time,
		}
		tagJSON, err := json.Marshal(e.Tags.Tags)
		if err != nil {
			log.Error().Err(err).Str("module", "storage").Msg("error while marshalling tag to JSON")
			return nil
		}
		b.Tags = tagJSON
		k, v = m.db.UpdateBucketToSQL(ctx, b)
	case resource.RESOURCE_TYPE_OBJECT:
		o := &models.Object{
			ObjectID:        e.Id.BigInt().String(),
			UpdateAt:        block.Block.Height,
			UpdateTxHash:    txHash,
			UpdateEVMTxHash: evmTxHash,
			UpdateTime:      block.Block.Time,
		}
		tagJSON, err := json.Marshal(e.Tags.Tags)
		if err != nil {
			log.Error().Err(err).Str("module", "storage").Msg("error while marshalling tag to JSON")
			return nil
		}
		o.Tags = tagJSON
		k, v = m.db.UpdateObjectToSQL(ctx, o)
	case resource.RESOURCE_TYPE_GROUP:
		g := &models.Group{
			GroupID:         e.Id.BigInt().String(),
			UpdateAt:        block.Block.Height,
			UpdateTxHash:    txHash,
			UpdateEVMTxHash: evmTxHash,
			UpdateTime:      block.Block.Time,
		}
		tagJSON, err := json.Marshal(e.Tags.Tags)
		if err != nil {
			log.Error().Err(err).Str("module", "storage").Msg("error while marshalling tag to JSON")
			return nil
		}
		g.Tags = tagJSON
		k, v = m.db.UpdateGroupToSQL(ctx, g)
	default:
		return nil
	}

	return map[string][]interface{}{
		k: v,
	}
}
