package storage

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	storagetypes "github.com/evmos/evmos/v12/x/storage/types"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
	"github.com/lib/pq"
)

var (
	EventCreateObject               = proto.MessageName(&storagetypes.EventCreateObject{})
	EventCancelCreateObject         = proto.MessageName(&storagetypes.EventCancelCreateObject{})
	EventSealObject                 = proto.MessageName(&storagetypes.EventSealObject{})
	EventCopyObject                 = proto.MessageName(&storagetypes.EventCopyObject{})
	EventDeleteObject               = proto.MessageName(&storagetypes.EventDeleteObject{})
	EventRejectSealObject           = proto.MessageName(&storagetypes.EventRejectSealObject{})
	EventDiscontinueObject          = proto.MessageName(&storagetypes.EventDiscontinueObject{})
	EventUpdateObjectInfo           = proto.MessageName(&storagetypes.EventUpdateObjectInfo{})
	EventUpdateObjectContent        = proto.MessageName(&storagetypes.EventUpdateObjectContent{})
	EventUpdateObjectContentSuccess = proto.MessageName(&storagetypes.EventUpdateObjectContentSuccess{})
	EventCancelUpdateObjectContent  = proto.MessageName(&storagetypes.EventCancelUpdateObjectContent{})
	EventMirrorObject               = proto.MessageName(&storagetypes.EventMirrorObject{})
	EventMirrorObjectResult         = proto.MessageName(&storagetypes.EventMirrorObjectResult{})
)

var ObjectEvents = map[string]bool{
	EventCreateObject:               true,
	EventCancelCreateObject:         true,
	EventSealObject:                 true,
	EventCopyObject:                 true,
	EventDeleteObject:               true,
	EventRejectSealObject:           true,
	EventDiscontinueObject:          true,
	EventUpdateObjectInfo:           true,
	EventUpdateObjectContent:        true,
	EventUpdateObjectContentSuccess: true,
	EventCancelUpdateObjectContent:  true,
	EventMirrorObject:               true,
	EventMirrorObjectResult:         true,
}

func (m *Module) ExtractObjectEventStatements(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, event sdk.Event) (interface{}, error) {
	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case EventCreateObject:
		createObject, ok := typedEvent.(*storagetypes.EventCreateObject)
		if !ok {
			return nil, errors.New("create object event assert error")
		}
		return m.handleCreateObject(ctx, block, txHash, evmTxHash, createObject), nil
	case EventCancelCreateObject:
		cancelCreateObject, ok := typedEvent.(*storagetypes.EventCancelCreateObject)
		if !ok {
			return nil, errors.New("cancel create object event assert error")
		}
		return m.handleCancelCreateObject(ctx, block, txHash, evmTxHash, cancelCreateObject), nil
	case EventSealObject:
		sealObject, ok := typedEvent.(*storagetypes.EventSealObject)
		if !ok {
			return nil, errors.New("seal object event assert error")
		}
		return m.handleSealObject(ctx, block, txHash, evmTxHash, sealObject), nil
	case EventCopyObject:
		copyObject, ok := typedEvent.(*storagetypes.EventCopyObject)
		if !ok {
			return nil, errors.New("copy object event assert error")
		}
		return m.handleCopyObject(ctx, block, txHash, evmTxHash, copyObject)
	case EventDeleteObject:
		deleteObject, ok := typedEvent.(*storagetypes.EventDeleteObject)
		if !ok {
			return nil, errors.New("delete object event assert error")
		}
		return m.handleDeleteObject(ctx, block, txHash, evmTxHash, deleteObject), nil
	case EventRejectSealObject:
		rejectSealObject, ok := typedEvent.(*storagetypes.EventRejectSealObject)
		if !ok {
			return nil, errors.New("reject seal object event assert error")
		}
		return m.handleRejectSealObject(ctx, block, txHash, evmTxHash, rejectSealObject), nil
	case EventDiscontinueObject:
		discontinueObject, ok := typedEvent.(*storagetypes.EventDiscontinueObject)
		if !ok {
			return nil, errors.New("discontinue object event assert error")
		}
		return m.handleEventDiscontinueObject(ctx, block, txHash, evmTxHash, discontinueObject), nil
	case EventUpdateObjectInfo:
		updateObjectInfo, ok := typedEvent.(*storagetypes.EventUpdateObjectInfo)
		if !ok {
			return nil, errors.New("update object event assert error")
		}
		return m.handleUpdateObjectInfo(ctx, block, txHash, evmTxHash, updateObjectInfo), nil
	case EventUpdateObjectContent:
		updateObjectContent, ok := typedEvent.(*storagetypes.EventUpdateObjectContent)
		if !ok {
			return nil, errors.New("update object event assert error")
		}
		return m.handleUpdateObjectContent(ctx, block, txHash, evmTxHash, updateObjectContent), nil
	case EventUpdateObjectContentSuccess:
		updateObjectContent, ok := typedEvent.(*storagetypes.EventUpdateObjectContentSuccess)
		if !ok {
			return nil, errors.New("update object success event assert error")
		}
		return m.handleUpdateObjectContentSuccess(ctx, block, txHash, evmTxHash, updateObjectContent), nil
	case EventCancelUpdateObjectContent:
		cancelUpdateObjectContent, ok := typedEvent.(*storagetypes.EventCancelUpdateObjectContent)
		if !ok {
			return nil, errors.New("cancel update object event assert error")
		}
		return m.handleCancelUpdateObjectContent(ctx, block, txHash, evmTxHash, cancelUpdateObjectContent), nil
	case EventMirrorObject, EventMirrorObjectResult:
		return nil, nil
	}
	return nil, nil
}

func ByteArrayToPqArray(ba [][]byte) pq.StringArray {
	pqArray := pq.StringArray{}
	for _, b := range ba {
		pqArray = append(pqArray, hex.EncodeToString(b))
	}
	return pqArray
}

func (m *Module) handleCreateObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, createObject *storagetypes.EventCreateObject) map[string][]interface{} {
	o := &models.Object{
		BucketID:        createObject.BucketId.BigInt().String(),
		BucketName:      createObject.BucketName,
		ObjectID:        createObject.ObjectId.BigInt().String(),
		ObjectName:      createObject.ObjectName,
		CreatorAddress:  createObject.Creator,
		OwnerAddress:    createObject.Owner,
		PayloadSize:     createObject.PayloadSize,
		Visibility:      createObject.Visibility.String(),
		ContentType:     createObject.ContentType,
		Status:          createObject.Status.String(),
		RedundancyType:  createObject.RedundancyType.String(),
		SourceType:      createObject.SourceType.String(),
		CheckSums:       ByteArrayToPqArray(createObject.Checksums),
		CreateTxHash:    txHash,
		CreateEVMTxHash: evmTxHash,
		CreateAt:        block.Block.Height,
		CreateTime:      time.Unix(createObject.CreateAt, 0),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      time.Unix(createObject.CreateAt, 0),
		Removed:         false,
	}
	k, v := m.db.SaveObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, o.ToObjectEvent(EventCreateObject))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleSealObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, sealObject *storagetypes.EventSealObject) map[string][]interface{} {
	o := &models.Object{
		BucketName:          sealObject.BucketName,
		ObjectName:          sealObject.ObjectName,
		ObjectID:            sealObject.ObjectId.BigInt().String(),
		OperatorAddress:     sealObject.Operator,
		LocalVirtualGroupID: sealObject.LocalVirtualGroupId,
		Status:              sealObject.Status.String(),
		SealedTxHash:        txHash,
		SealedEvmTxHash:     evmTxHash,
		CheckSums:           ByteArrayToPqArray(sealObject.GetChecksums()),
		UpdateAt:            block.Block.Height,
		UpdateTxHash:        txHash,
		UpdateTime:          block.Block.Time,
		Removed:             false,
	}

	res := make(map[string][]interface{})

	k, v := m.db.UpdateStorageSizeToSQL(ctx, sealObject.ObjectId.BigInt().String(), sealObject.BucketName, "+")
	res[k] = v

	k, v = m.db.UpdateChargeSizeToSQL(ctx, sealObject.ObjectId.BigInt().String(), sealObject.BucketName, "+")
	res[k] = v

	k, v = m.db.UpdateObjectToSQL(ctx, o)
	res[k] = v

	k, v = m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventSealObject))
	res[k] = v

	return res
}

func (m *Module) handleCancelCreateObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, cancelCreateObject *storagetypes.EventCancelCreateObject) map[string][]interface{} {
	o := &models.Object{
		BucketName:      cancelCreateObject.BucketName,
		ObjectName:      cancelCreateObject.ObjectName,
		ObjectID:        cancelCreateObject.ObjectId.BigInt().String(),
		OperatorAddress: cancelCreateObject.Operator,
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
		Removed:         true,
	}

	k, v := m.db.UpdateObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventCancelCreateObject))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleCopyObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, copyObject *storagetypes.EventCopyObject) (map[string][]interface{}, error) {
	o, err := m.db.GetObject(ctx, copyObject.SrcObjectId.BigInt().String())
	if err != nil {
		return nil, err
	}
	o.ObjectID = copyObject.DstObjectId.BigInt().String()
	o.ObjectName = copyObject.DstObjectName
	o.BucketName = copyObject.DstBucketName
	o.OperatorAddress = copyObject.Operator
	o.CreateAt = block.Block.Height
	o.CreateTxHash = txHash
	o.CreateTime = block.Block.Time
	o.UpdateAt = block.Block.Height
	o.UpdateTxHash = txHash
	o.UpdateEVMTxHash = evmTxHash
	o.UpdateTime = block.Block.Time
	o.Removed = false
	if o.PayloadSize == 0 {
		o.Status = storagetypes.OBJECT_STATUS_SEALED.String()
	} else {
		o.Status = storagetypes.OBJECT_STATUS_CREATED.String()
	}

	res := make(map[string][]interface{})

	k, v := m.db.SaveObjectToSQL(ctx, o)
	res[k] = v

	if o.PayloadSize == 0 {
		k, v = m.db.UpdateChargeSizeToSQL(ctx, copyObject.DstObjectId.BigInt().String(), copyObject.DstBucketName, "+")
		res[k] = v
	}

	k, v = m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventCopyObject))
	res[k] = v

	return res, nil
}

func (m *Module) handleDeleteObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, deleteObject *storagetypes.EventDeleteObject) map[string][]interface{} {
	o := &models.Object{
		BucketName:          deleteObject.BucketName,
		ObjectName:          deleteObject.ObjectName,
		ObjectID:            deleteObject.ObjectId.BigInt().String(),
		LocalVirtualGroupID: deleteObject.LocalVirtualGroupId,
		UpdateAt:            block.Block.Height,
		UpdateTxHash:        txHash,
		UpdateEVMTxHash:     evmTxHash,
		UpdateTime:          block.Block.Time,
		Removed:             true,
	}

	res := make(map[string][]interface{})
	k, v := m.db.UpdateStorageSizeToSQL(ctx, deleteObject.ObjectId.BigInt().String(), deleteObject.BucketName, "-")
	res[k] = v

	k, v = m.db.UpdateChargeSizeToSQL(ctx, deleteObject.ObjectId.BigInt().String(), deleteObject.BucketName, "-")
	res[k] = v

	k, v = m.db.UpdateObjectToSQL(ctx, o)
	res[k] = v

	k, v = m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventDeleteObject))
	res[k] = v

	return res
}

// RejectSeal event won't emit a delete event, need to be deleted manually here in metadata service
// handle logic is set as removed, no need to set status
func (m *Module) handleRejectSealObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, rejectSealObject *storagetypes.EventRejectSealObject) map[string][]interface{} {
	o := &models.Object{
		BucketName:      rejectSealObject.BucketName,
		ObjectName:      rejectSealObject.ObjectName,
		ObjectID:        rejectSealObject.ObjectId.BigInt().String(),
		OperatorAddress: rejectSealObject.Operator,
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
	}
	if rejectSealObject.ForUpdate {
		o.IsUpdating = false
	} else {
		o.Removed = true
	}
	k, v := m.db.UpdateObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventRejectSealObject))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleEventDiscontinueObject(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, discontinueObject *storagetypes.EventDiscontinueObject) map[string][]interface{} {
	o := &models.Object{
		BucketName:      discontinueObject.BucketName,
		ObjectID:        discontinueObject.ObjectId.BigInt().String(),
		DeleteReason:    discontinueObject.Reason,
		DeleteAt:        discontinueObject.DeleteAt,
		Status:          storagetypes.OBJECT_STATUS_DISCONTINUED.String(),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
		Removed:         false,
	}
	k, v := m.db.UpdateObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventDiscontinueObject))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleUpdateObjectInfo(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, updateObject *storagetypes.EventUpdateObjectInfo) map[string][]interface{} {
	o := &models.Object{
		BucketName:      updateObject.BucketName,
		ObjectID:        updateObject.ObjectId.BigInt().String(),
		ObjectName:      updateObject.ObjectName,
		OperatorAddress: updateObject.Operator,
		Visibility:      updateObject.Visibility.String(),
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
	}
	k, v := m.db.UpdateObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventUpdateObjectInfo))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleUpdateObjectContent(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, updateObject *storagetypes.EventUpdateObjectContent) map[string][]interface{} {
	o := &models.Object{
		BucketName:      updateObject.BucketName,
		ObjectID:        updateObject.ObjectId.BigInt().String(),
		ObjectName:      updateObject.ObjectName,
		OperatorAddress: updateObject.Operator,
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,

		IsUpdating: true,
	}
	k, v := m.db.UpdateObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, models.NewObjectEvent(o.ObjectID, block.Block.Height, txHash, evmTxHash, EventUpdateObjectContent))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

// handleUpdateObjectContentSuccess, when sealing an updated object, EventUpdateObjectContentSuccess will be emitted before EventSealObjet.
// Returns ordered []database.SQLStatement so PostgreSQL receives one statement per Exec (no multi-statement concatenation, MOCA26Q1-003).
func (m *Module) handleUpdateObjectContentSuccess(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, updateObject *storagetypes.EventUpdateObjectContentSuccess) []database.SQLStatement {
	o := &models.Object{
		BucketName:         updateObject.BucketName,
		ObjectID:           updateObject.ObjectId.BigInt().String(),
		ObjectName:         updateObject.ObjectName,
		OperatorAddress:    updateObject.Operator,
		UpdateAt:           block.Block.Height,
		UpdateTxHash:       txHash,
		UpdateEVMTxHash:    evmTxHash,
		UpdateTime:         block.Block.Time,
		ContentType:        updateObject.ContentType,
		IsUpdating:         false,
		ContentUpdatedTime: time.Unix(updateObject.UpdatedAt, 0),
		Updater:            updateObject.Operator,
		PayloadSize:        updateObject.NewPayloadSize,
		CheckSums:          ByteArrayToPqArray(updateObject.NewChecksums),
		Version:            updateObject.Version,
	}
	k1, v1 := m.db.UpdateStorageSizeToSQL(ctx, updateObject.ObjectId.BigInt().String(), updateObject.BucketName, "-")
	k2, v2 := m.db.UpdateChargeSizeToSQL(ctx, updateObject.ObjectId.BigInt().String(), updateObject.BucketName, "-")
	k3, v3 := m.db.UpdateObjectToSQL(ctx, o)
	k4, v4 := m.db.SaveObjectEventToSQL(ctx, o.ToObjectEvent(EventUpdateObjectContentSuccess))
	return []database.SQLStatement{
		{SQL: k1, Vars: v1},
		{SQL: k2, Vars: v2},
		{SQL: k3, Vars: v3},
		{SQL: k4, Vars: v4},
	}
}

func (m *Module) handleCancelUpdateObjectContent(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, cancelUpdateObject *storagetypes.EventCancelUpdateObjectContent) map[string][]interface{} {
	o := &models.Object{
		BucketName:      cancelUpdateObject.BucketName,
		ObjectName:      cancelUpdateObject.ObjectName,
		ObjectID:        cancelUpdateObject.ObjectId.BigInt().String(),
		OperatorAddress: cancelUpdateObject.Operator,
		UpdateAt:        block.Block.Height,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		UpdateTime:      block.Block.Time,
		IsUpdating:      false,
	}
	k, v := m.db.UpdateObjectToSQL(ctx, o)
	ek, ev := m.db.SaveObjectEventToSQL(ctx, o.ToObjectEvent(EventCancelUpdateObjectContent))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}
