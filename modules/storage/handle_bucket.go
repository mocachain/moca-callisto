package storage

import (
	"context"
	"errors"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	storagetypes "github.com/evmos/evmos/v12/x/storage/types"
	"github.com/rs/zerolog/log"
)

var (
	EventCreateBucket             = proto.MessageName(&storagetypes.EventCreateBucket{})
	EventDeleteBucket             = proto.MessageName(&storagetypes.EventDeleteBucket{})
	EventUpdateBucketInfo         = proto.MessageName(&storagetypes.EventUpdateBucketInfo{})
	EventDiscontinueBucket        = proto.MessageName(&storagetypes.EventDiscontinueBucket{})
	EventMigrationBucket          = proto.MessageName(&storagetypes.EventMigrationBucket{})
	EventCancelMigrationBucket    = proto.MessageName(&storagetypes.EventCancelMigrationBucket{})
	EventRejectMigrateBucket      = proto.MessageName(&storagetypes.EventRejectMigrateBucket{})
	EventCompleteMigrationBucket  = proto.MessageName(&storagetypes.EventCompleteMigrationBucket{})
	EventToggleSPAsDelegatedAgent = proto.MessageName(&storagetypes.EventToggleSPAsDelegatedAgent{})
	EventMirrorBucket             = proto.MessageName(&storagetypes.EventMirrorBucket{})
	EventMirrorBucketResult       = proto.MessageName(&storagetypes.EventMirrorBucketResult{})
)

var BucketEvents = map[string]bool{
	EventCreateBucket:             true,
	EventDeleteBucket:             true,
	EventUpdateBucketInfo:         true,
	EventDiscontinueBucket:        true,
	EventMigrationBucket:          true,
	EventCancelMigrationBucket:    true,
	EventRejectMigrateBucket:      true,
	EventCompleteMigrationBucket:  true,
	EventToggleSPAsDelegatedAgent: true,
	EventMirrorBucket:             true,
	EventMirrorBucketResult:       true,
}

func (m *Module) ExtractBucketEventStatements(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, event sdk.Event) (interface{}, error) {
	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case EventCreateBucket:
		createBucket, ok := typedEvent.(*storagetypes.EventCreateBucket)
		if !ok {
			return nil, errors.New("create bucket event assert error")
		}
		return m.handleCreateBucket(ctx, block, txHash, evmTxHash, createBucket), nil
	case EventDeleteBucket:
		deleteBucket, ok := typedEvent.(*storagetypes.EventDeleteBucket)
		if !ok {
			return nil, errors.New("delete bucket event assert error")
		}
		return m.handleDeleteBucket(ctx, block, txHash, evmTxHash, deleteBucket), nil
	case EventUpdateBucketInfo:
		updateBucketInfo, ok := typedEvent.(*storagetypes.EventUpdateBucketInfo)
		if !ok {
			return nil, errors.New("update bucket event assert error")
		}
		return m.handleUpdateBucketInfo(ctx, block, txHash, evmTxHash, updateBucketInfo), nil
	case EventDiscontinueBucket:
		discontinueBucket, ok := typedEvent.(*storagetypes.EventDiscontinueBucket)
		if !ok {
			return nil, errors.New("discontinue bucket event assert error")
		}
		return m.handleDiscontinueBucket(ctx, block, txHash, evmTxHash, discontinueBucket), nil
	case EventMigrationBucket:
		migrationBucket, ok := typedEvent.(*storagetypes.EventMigrationBucket)
		if !ok {
			return nil, errors.New("migration bucket event assert error")
		}
		return m.handleEventMigrationBucket(ctx, block, txHash, evmTxHash, migrationBucket), nil
	case EventCancelMigrationBucket:
		cancelMigrationBucket, ok := typedEvent.(*storagetypes.EventCancelMigrationBucket)
		if !ok {
			return nil, errors.New("cancel migration bucket event assert error")
		}
		return m.handleEventCancelMigrationBucket(ctx, block, txHash, evmTxHash, cancelMigrationBucket), nil
	case EventRejectMigrateBucket:
		rejectMigrateBucket, ok := typedEvent.(*storagetypes.EventRejectMigrateBucket)
		if !ok {
			return nil, errors.New("reject migration bucket event assert error")
		}
		return m.handleEventRejectMigrateBucket(ctx, block, txHash, evmTxHash, rejectMigrateBucket), nil
	case EventCompleteMigrationBucket:
		completeMigrationBucket, ok := typedEvent.(*storagetypes.EventCompleteMigrationBucket)
		if !ok {
			return nil, errors.New("complete migrate bucket event assert error")
		}
		return m.handleCompleteMigrationBucket(ctx, block, txHash, evmTxHash, completeMigrationBucket), nil
	case EventToggleSPAsDelegatedAgent:
		toggleSP, ok := typedEvent.(*storagetypes.EventToggleSPAsDelegatedAgent)
		if !ok {
			return nil, errors.New("toggle sp as delegated agent event assert error")
		}
		return m.handleToggleSPAsDelegatedAgent(ctx, block, txHash, evmTxHash, toggleSP), nil
	case EventMirrorBucket, EventMirrorBucketResult:
		return nil, nil
	}

	return nil, nil
}

func (m *Module) handleCreateBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, createBucket *storagetypes.EventCreateBucket) map[string][]interface{} {
	gvgFamilies, err := m.vgSource.GlobalVirtualGroupFamilies(block.Block.Height, context.Background())
	if err != nil {
		return nil
	}

	for _, gvgf := range gvgFamilies {
		s := &models.GlobalVirtualGroupFamily{
			GlobalVirtualGroupFamilyID: gvgf.Id,
			PrimarySpID:                gvgf.PrimarySpId,
			GlobalVirtualGroupIDs:      models.ConvertUint32ToInt32Array(gvgf.GlobalVirtualGroupIds),
			VirtualPaymentAddress:      gvgf.VirtualPaymentAddress,
			CreateAt:                   1,
		}
		k, v := m.db.SaveGVGFToSQL(ctx, s)
		if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
			log.Err(err).Msg("failed to save global virtual group family")
			continue
		}
		gvgs, err := m.vgSource.GlobalVirtualGroupByFamilyID(block.Block.Height, gvgf.Id)
		if err != nil {
			log.Err(err).Msg("failed to get global virtual group by family id")
			continue
		}
		for _, gvg := range gvgs {
			g := &models.GlobalVirtualGroup{
				GlobalVirtualGroupID:  gvg.Id,
				FamilyID:              gvg.FamilyId,
				PrimarySpID:           gvg.PrimarySpId,
				SecondarySpIDs:        models.ConvertUint32ToInt32Array(gvg.SecondarySpIds),
				StoredSize:            gvg.StoredSize,
				VirtualPaymentAddress: gvg.VirtualPaymentAddress,
				TotalDeposit:          gvg.TotalDeposit.BigInt().String(),
			}
			k, v := m.db.SaveGVGToSQL(ctx, g)
			if err := m.db.ExecuteStatements([]database.SQLStatement{{SQL: k, Vars: v}}); err != nil {
				log.Err(err).Msg("failed to save global virtual group")
				continue
			}
		}
	}

	b := &models.Bucket{
		BucketID:                   createBucket.BucketId.BigInt().String(),
		BucketName:                 createBucket.BucketName,
		OwnerAddress:               createBucket.Owner,
		PaymentAddress:             createBucket.PaymentAddress,
		GlobalVirtualGroupFamilyId: createBucket.GlobalVirtualGroupFamilyId,
		OperatorAddress:            createBucket.Owner,
		SourceType:                 createBucket.SourceType.String(),
		ChargedReadQuota:           createBucket.ChargedReadQuota,
		Visibility:                 createBucket.Visibility.String(),
		Status:                     createBucket.Status.String(),
		Removed:                    false,
		CreateAt:                   block.Block.Height,
		CreateTxHash:               txHash,
		CreateEVMTxHash:            evmTxHash,
		CreateTime:                 time.Unix(createBucket.CreateAt, 0),
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateEVMTxHash:            evmTxHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.SaveBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, b.ToBucketEvent(EventCreateBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleDeleteBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, deleteBucket *storagetypes.EventDeleteBucket) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:                   deleteBucket.BucketId.BigInt().String(),
		BucketName:                 deleteBucket.BucketName,
		OwnerAddress:               deleteBucket.Owner,
		GlobalVirtualGroupFamilyId: deleteBucket.GlobalVirtualGroupFamilyId,
		DeleteAt:                   block.Block.Height,
		Removed:                    true,
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateEVMTxHash:            evmTxHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, models.NewBucketEvent(b.BucketID, block.Block.Height, txHash, evmTxHash, EventDeleteBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleDiscontinueBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, discontinueBucket *storagetypes.EventDiscontinueBucket) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:        discontinueBucket.BucketId.BigInt().String(),
		BucketName:      discontinueBucket.BucketName,
		DeleteReason:    discontinueBucket.Reason,
		DeleteAt:        discontinueBucket.DeleteAt,
		Status:          storagetypes.BUCKET_STATUS_DISCONTINUED.String(),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
	}
	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, models.NewBucketEvent(b.BucketID, block.Block.Height, txHash, evmTxHash, EventDiscontinueBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleUpdateBucketInfo(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, updateBucket *storagetypes.EventUpdateBucketInfo) map[string][]interface{} {
	b := &models.Bucket{
		BucketName:                 updateBucket.BucketName,
		BucketID:                   updateBucket.BucketId.BigInt().String(),
		ChargedReadQuota:           updateBucket.ChargedReadQuota,
		PaymentAddress:             updateBucket.PaymentAddress,
		Visibility:                 updateBucket.Visibility.String(),
		GlobalVirtualGroupFamilyId: updateBucket.GlobalVirtualGroupFamilyId,
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateEVMTxHash:            evmTxHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, models.NewBucketEvent(b.BucketID, block.Block.Height, txHash, evmTxHash, EventUpdateBucketInfo))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleEventMigrationBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, migrationBucket *storagetypes.EventMigrationBucket) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:        migrationBucket.BucketId.BigInt().String(),
		BucketName:      migrationBucket.BucketName,
		Status:          migrationBucket.Status.String(),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
	}

	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, b.ToBucketEvent(EventMigrationBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleEventCancelMigrationBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, cancelMigrationBucket *storagetypes.EventCancelMigrationBucket) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:        cancelMigrationBucket.BucketId.BigInt().String(),
		BucketName:      cancelMigrationBucket.BucketName,
		Status:          cancelMigrationBucket.Status.String(),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
	}
	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, b.ToBucketEvent(EventCancelMigrationBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleEventRejectMigrateBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, rejectMigrateBucket *storagetypes.EventRejectMigrateBucket) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:        rejectMigrateBucket.BucketId.BigInt().String(),
		BucketName:      rejectMigrateBucket.BucketName,
		Status:          rejectMigrateBucket.Status.String(),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
	}
	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, b.ToBucketEvent(EventRejectMigrateBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleCompleteMigrationBucket(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, completeMigrationBucket *storagetypes.EventCompleteMigrationBucket) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:                   completeMigrationBucket.BucketId.BigInt().String(),
		BucketName:                 completeMigrationBucket.BucketName,
		GlobalVirtualGroupFamilyId: completeMigrationBucket.GlobalVirtualGroupFamilyId,
		Status:                     completeMigrationBucket.Status.String(),
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateEVMTxHash:            evmTxHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.UpdateBucketToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, b.ToBucketEvent(EventCompleteMigrationBucket))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleToggleSPAsDelegatedAgent(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, toggleSP *storagetypes.EventToggleSPAsDelegatedAgent) map[string][]interface{} {
	b := &models.Bucket{
		BucketID:                   toggleSP.BucketId.BigInt().String(),
		SpAsDelegatedAgentDisabled: toggleSP.SpAsDelegatedAgentDisabled,
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateEVMTxHash:            evmTxHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.UpdateBucketSPAsDelegatedAgentToSQL(ctx, b)
	ek, ev := m.db.SaveBucketEventToSQL(ctx, models.NewBucketEvent(b.BucketID, block.Block.Height, txHash, evmTxHash, EventToggleSPAsDelegatedAgent))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}
