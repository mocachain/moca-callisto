package virtualgroup

import (
	"context"
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	vgtypes "github.com/evmos/evmos/v12/x/virtualgroup/types"
	junotypes "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
)

var (
	EventCreateLocalVirtualGroup        = proto.MessageName(&vgtypes.EventCreateLocalVirtualGroup{})
	EventDeleteLocalVirtualGroup        = proto.MessageName(&vgtypes.EventDeleteLocalVirtualGroup{})
	EventUpdateLocalVirtualGroup        = proto.MessageName(&vgtypes.EventUpdateLocalVirtualGroup{})
	EventCreateGlobalVirtualGroup       = proto.MessageName(&vgtypes.EventCreateGlobalVirtualGroup{})
	EventDeleteGlobalVirtualGroup       = proto.MessageName(&vgtypes.EventDeleteGlobalVirtualGroup{})
	EventUpdateGlobalVirtualGroup       = proto.MessageName(&vgtypes.EventUpdateGlobalVirtualGroup{})
	EventCreateGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventCreateGlobalVirtualGroupFamily{})
	EventDeleteGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventDeleteGlobalVirtualGroupFamily{})
	EventUpdateGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventUpdateGlobalVirtualGroupFamily{})
)

var VirtualGroupEvents = map[string]bool{
	EventCreateLocalVirtualGroup:        true,
	EventDeleteLocalVirtualGroup:        true,
	EventUpdateLocalVirtualGroup:        true,
	EventCreateGlobalVirtualGroup:       true,
	EventDeleteGlobalVirtualGroup:       true,
	EventUpdateGlobalVirtualGroup:       true,
	EventCreateGlobalVirtualGroupFamily: true,
	EventDeleteGlobalVirtualGroupFamily: true,
	EventUpdateGlobalVirtualGroupFamily: true,
}

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
	allSQL := make([]database.SQLStatement, 0)
	for _, event := range tx.Events {
		sqls, err := m.ExtractEventStatements(ctx, block, tx.TxHash, sdk.Event(event))
		if err != nil {
			log.Err(err)
			continue
		}
		for k, v := range sqls {
			allSQL = append(allSQL, database.SQLStatement{SQL: k, Vars: v})
		}
	}
	return allSQL, nil
}

func (m *Module) ExtractEventStatements(ctx context.Context, block *tmctypes.ResultBlock, txHash string, event sdk.Event) (map[string][]interface{}, error) {
	if !VirtualGroupEvents[event.Type] {
		return nil, nil
	}
	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case EventCreateLocalVirtualGroup:
		createLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventCreateLocalVirtualGroup)
		if !ok {
			return nil, errors.New("create lvg event assert error")
		}
		return m.handleCreateLocalVirtualGroup(ctx, block, txHash, createLocalVirtualGroup), nil
	case EventDeleteLocalVirtualGroup:
		deleteLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventDeleteLocalVirtualGroup)
		if !ok {
			return nil, errors.New("delete lvg event assert error")
		}
		return m.handleDeleteLocalVirtualGroup(ctx, block, txHash, deleteLocalVirtualGroup), nil
	case EventUpdateLocalVirtualGroup:
		updateLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventUpdateLocalVirtualGroup)
		if !ok {
			return nil, errors.New("update lvg event assert error")
		}
		return m.handleUpdateLocalVirtualGroup(ctx, block, txHash, updateLocalVirtualGroup), nil
	case EventCreateGlobalVirtualGroup:
		createGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventCreateGlobalVirtualGroup)
		if !ok {
			return nil, errors.New("create gvg event assert error")
		}
		return m.handleCreateGlobalVirtualGroup(ctx, block, txHash, createGlobalVirtualGroup), nil
	case EventDeleteGlobalVirtualGroup:
		deleteGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventDeleteGlobalVirtualGroup)
		if !ok {
			return nil, errors.New("delete gvg event assert error")
		}
		return m.handleDeleteGlobalVirtualGroup(ctx, block, txHash, deleteGlobalVirtualGroup), nil
	case EventUpdateGlobalVirtualGroup:
		updateGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventUpdateGlobalVirtualGroup)
		if !ok {
			return nil, errors.New("update gvg event assert error")
		}
		return m.handleUpdateGlobalVirtualGroup(ctx, block, txHash, updateGlobalVirtualGroup), nil
	case EventCreateGlobalVirtualGroupFamily:
		createGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventCreateGlobalVirtualGroupFamily)
		if !ok {
			return nil, errors.New("create vgf event assert error")
		}
		return m.handleCreateGlobalVirtualGroupFamily(ctx, block, txHash, createGlobalVirtualGroupFamily), nil
	case EventDeleteGlobalVirtualGroupFamily:
		deleteGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventDeleteGlobalVirtualGroupFamily)
		if !ok {
			return nil, errors.New("delete vgf event assert error")
		}
		return m.handleDeleteGlobalVirtualGroupFamily(ctx, block, txHash, deleteGlobalVirtualGroupFamily), nil
	case EventUpdateGlobalVirtualGroupFamily:
		updateGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventUpdateGlobalVirtualGroupFamily)
		if !ok {
			return nil, errors.New("update vgf event assert error")
		}
		return m.handleUpdateGlobalVirtualGroupFamily(ctx, block, txHash, updateGlobalVirtualGroupFamily), nil
	}
	return nil, nil
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash string, event sdk.Event) error {
	return nil
}

func (m *Module) handleCreateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash string, createLocalVirtualGroup *vgtypes.EventCreateLocalVirtualGroup) map[string][]interface{} {
	lvgGroup := &models.LocalVirtualGroup{
		LocalVirtualGroupID:  createLocalVirtualGroup.Id,
		GlobalVirtualGroupID: createLocalVirtualGroup.GlobalVirtualGroupId,
		BucketID:             createLocalVirtualGroup.BucketId.BigInt().String(),
		StoredSize:           createLocalVirtualGroup.StoredSize,
		CreateAt:             block.Block.Height,
		CreateTxHash:         txHash,
		CreateTime:           block.Block.Time,
		UpdateAt:             block.Block.Height,
		UpdateTxHash:         txHash,
		UpdateTime:           block.Block.Time,
		Removed:              false,
	}
	k, v := m.db.SaveLVGToSQL(ctx, lvgGroup)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleUpdateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash string, updateLocalVirtualGroup *vgtypes.EventUpdateLocalVirtualGroup) map[string][]interface{} {
	lvgGroup := &models.LocalVirtualGroup{
		LocalVirtualGroupID:  updateLocalVirtualGroup.Id,
		BucketID:             updateLocalVirtualGroup.BucketId.BigInt().String(),
		GlobalVirtualGroupID: updateLocalVirtualGroup.GlobalVirtualGroupId,
		StoredSize:           updateLocalVirtualGroup.StoredSize,
		UpdateAt:             block.Block.Height,
		UpdateTxHash:         txHash,
		UpdateTime:           block.Block.Time,
	}
	k, v := m.db.UpdateLVGToSQL(ctx, lvgGroup)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleCreateGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash string, createGlobalVirtualGroup *vgtypes.EventCreateGlobalVirtualGroup) map[string][]interface{} {
	gvgGroup := &models.GlobalVirtualGroup{
		GlobalVirtualGroupID:  createGlobalVirtualGroup.Id,
		FamilyID:              createGlobalVirtualGroup.FamilyId,
		PrimarySpID:           createGlobalVirtualGroup.PrimarySpId,
		SecondarySpIDs:        models.ConvertUint32ToInt32Array(createGlobalVirtualGroup.SecondarySpIds),
		StoredSize:            createGlobalVirtualGroup.StoredSize,
		VirtualPaymentAddress: createGlobalVirtualGroup.VirtualPaymentAddress,
		TotalDeposit:          createGlobalVirtualGroup.TotalDeposit.BigInt().String(),
		CreateAt:              block.Block.Height,
		CreateTxHash:          txHash,
		CreateTime:            block.Block.Time,
		UpdateAt:              block.Block.Height,
		UpdateTxHash:          txHash,
		UpdateTime:            block.Block.Time,
		Removed:               false,
	}

	k, v := m.db.SaveGVGToSQL(ctx, gvgGroup)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleDeleteGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash string, deleteGlobalVirtualGroup *vgtypes.EventDeleteGlobalVirtualGroup) map[string][]interface{} {
	gvgGroup := &models.GlobalVirtualGroup{
		GlobalVirtualGroupID: deleteGlobalVirtualGroup.Id,
		Removed:              true,
		UpdateAt:             block.Block.Height,
		UpdateTxHash:         txHash,
		UpdateTime:           block.Block.Time,
	}
	k, v := m.db.DeleteGVGToSQL(ctx, gvgGroup)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleUpdateGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash string, updateGlobalVirtualGroup *vgtypes.EventUpdateGlobalVirtualGroup) map[string][]interface{} {
	gvgGroup := &models.GlobalVirtualGroup{
		GlobalVirtualGroupID: updateGlobalVirtualGroup.Id,
		StoredSize:           updateGlobalVirtualGroup.StoreSize,
		TotalDeposit:         updateGlobalVirtualGroup.TotalDeposit.BigInt().String(),
		PrimarySpID:          updateGlobalVirtualGroup.PrimarySpId,
		SecondarySpIDs:       models.ConvertUint32ToInt32Array(updateGlobalVirtualGroup.SecondarySpIds),
		UpdateAt:             block.Block.Height,
		UpdateTxHash:         txHash,
		UpdateTime:           block.Block.Time,
	}

	k, v := m.db.UpdateGVGToSQL(ctx, gvgGroup)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleCreateGlobalVirtualGroupFamily(ctx context.Context, block *tmctypes.ResultBlock, txHash string, createGlobalVirtualGroupFamily *vgtypes.EventCreateGlobalVirtualGroupFamily) map[string][]interface{} {
	vgfGroup := &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyID: createGlobalVirtualGroupFamily.Id,
		PrimarySpID:                createGlobalVirtualGroupFamily.PrimarySpId,
		VirtualPaymentAddress:      createGlobalVirtualGroupFamily.VirtualPaymentAddress,
		GlobalVirtualGroupIDs:      models.ConvertUint32ToInt32Array(createGlobalVirtualGroupFamily.GlobalVirtualGroupIds),
		CreateAt:                   block.Block.Height,
		CreateTxHash:               txHash,
		CreateTime:                 block.Block.Time,
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateTime:                 block.Block.Time,
		Removed:                    false,
	}
	k, v := m.db.SaveGVGFToSQL(ctx, vgfGroup)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleDeleteLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash string, deleteLocalVirtualGroup *vgtypes.EventDeleteLocalVirtualGroup) map[string][]interface{} {
	vg := &models.LocalVirtualGroup{
		LocalVirtualGroupID: deleteLocalVirtualGroup.Id,
		BucketID:            deleteLocalVirtualGroup.BucketId.BigInt().String(),
		Removed:             true,
		UpdateAt:            block.Block.Height,
		UpdateTxHash:        txHash,
		UpdateTime:          block.Block.Time,
	}
	k, v := m.db.DeleteLVGToSQL(ctx, vg)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleDeleteGlobalVirtualGroupFamily(ctx context.Context, block *tmctypes.ResultBlock, txHash string, deleteGlobalVirtualGroupFamily *vgtypes.EventDeleteGlobalVirtualGroupFamily) map[string][]interface{} {
	vg := &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyID: deleteGlobalVirtualGroupFamily.Id,
		Removed:                    true,
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.UpdateVGFToSQL(ctx, vg)
	return map[string][]interface{}{
		k: v,
	}
}

func (m *Module) handleUpdateGlobalVirtualGroupFamily(ctx context.Context, block *tmctypes.ResultBlock, txHash string, updateGlobalVirtualGroupFamily *vgtypes.EventUpdateGlobalVirtualGroupFamily) map[string][]interface{} {
	vg := &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyID: updateGlobalVirtualGroupFamily.Id,
		PrimarySpID:                updateGlobalVirtualGroupFamily.PrimarySpId,
		GlobalVirtualGroupIDs:      models.ConvertUint32ToInt32Array(updateGlobalVirtualGroupFamily.GlobalVirtualGroupIds),
		UpdateAt:                   block.Block.Height,
		UpdateTxHash:               txHash,
		UpdateTime:                 block.Block.Time,
	}
	k, v := m.db.UpdateVGFToSQL(ctx, vg)
	return map[string][]interface{}{
		k: v,
	}
}
