package storage

import (
	"context"
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	storagetypes "github.com/evmos/evmos/v12/x/storage/types"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/models"
)

var (
	EventCreateGroup       = proto.MessageName(&storagetypes.EventCreateGroup{})
	EventDeleteGroup       = proto.MessageName(&storagetypes.EventDeleteGroup{})
	EventLeaveGroup        = proto.MessageName(&storagetypes.EventLeaveGroup{})
	EventUpdateGroupMember = proto.MessageName(&storagetypes.EventUpdateGroupMember{})
	EventRenewGroupMember  = proto.MessageName(&storagetypes.EventRenewGroupMember{})
	EventUpdateGroupExtra  = proto.MessageName(&storagetypes.EventUpdateGroupExtra{})
	EventMirrorGroup       = proto.MessageName(&storagetypes.EventMirrorGroup{})
	EventMirrorGroupResult = proto.MessageName(&storagetypes.EventMirrorGroupResult{})
)

var GroupEvents = map[string]bool{
	EventCreateGroup:       true,
	EventDeleteGroup:       true,
	EventLeaveGroup:        true,
	EventUpdateGroupMember: true,
	EventRenewGroupMember:  true,
	EventUpdateGroupExtra:  true,
	EventMirrorGroup:       true,
	EventMirrorGroupResult: true,
}

func (m *Module) ExtractGroupEventStatements(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, event sdk.Event) (interface{}, error) {
	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case EventCreateGroup:
		createGroup, ok := typedEvent.(*storagetypes.EventCreateGroup)
		if !ok {
			return nil, errors.New("create group event assert error")
		}
		return m.handleCreateGroup(ctx, block, txHash, evmTxHash, createGroup), nil
	case EventUpdateGroupMember:
		updateGroupMember, ok := typedEvent.(*storagetypes.EventUpdateGroupMember)
		if !ok {
			return nil, errors.New("update group member event assert error")
		}
		return m.handleUpdateGroupMember(ctx, block, txHash, evmTxHash, updateGroupMember), nil
	case EventDeleteGroup:
		deleteGroup, ok := typedEvent.(*storagetypes.EventDeleteGroup)
		if !ok {
			return nil, errors.New("delete group event assert error")
		}
		return m.handleDeleteGroup(ctx, block, txHash, evmTxHash, deleteGroup), nil
	case EventLeaveGroup:
		leaveGroup, ok := typedEvent.(*storagetypes.EventLeaveGroup)
		if !ok {
			return nil, errors.New("leave group event assert error")
		}
		return m.handleLeaveGroup(ctx, block, txHash, evmTxHash, leaveGroup), nil
	case EventRenewGroupMember:
		renewGroupMember, ok := typedEvent.(*storagetypes.EventRenewGroupMember)
		if !ok {
			return nil, errors.New("renew group member event assert error")
		}
		return m.handleRenewGroupMember(ctx, block, txHash, evmTxHash, renewGroupMember), nil
	case EventUpdateGroupExtra, EventMirrorGroup, EventMirrorGroupResult:
		return nil, nil
	}
	return nil, nil
}

func (m *Module) handleCreateGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, createGroup *storagetypes.EventCreateGroup) map[string][]interface{} {
	var membersToAddList []*models.Group
	g := &models.Group{
		OwnerAddress:    createGroup.Owner,
		GroupID:         createGroup.GroupId.BigInt().String(),
		GroupName:       createGroup.GroupName,
		SourceType:      createGroup.SourceType.String(),
		Extra:           createGroup.Extra,
		CreateAt:        block.Block.Height,
		CreateTxHash:    txHash,
		CreateEVMTxHash: evmTxHash,
		CreateTime:      block.Block.Time,
		UpdateAt:        block.Block.Height,
		UpdateTime:      block.Block.Time,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		Removed:         false,
	}
	membersToAddList = append(membersToAddList, g)
	k, v := m.db.CreateGroupToSQL(ctx, membersToAddList)
	ek, ev := m.db.SaveGroupEventToSQL(ctx, g.ToBucketEvent(EventCreateGroup))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleDeleteGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, deleteGroup *storagetypes.EventDeleteGroup) map[string][]interface{} {
	g := &models.Group{
		OwnerAddress:    deleteGroup.Owner,
		GroupID:         deleteGroup.GroupId.BigInt().String(),
		GroupName:       deleteGroup.GroupName,
		UpdateAt:        block.Block.Height,
		UpdateTime:      block.Block.Time,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		Removed:         true,
	}
	k, v := m.db.DeleteGroupToSQL(ctx, g)
	ek, ev := m.db.SaveGroupEventToSQL(ctx, models.NewGroupEvent(g.GroupID, block.Block.Height, txHash, evmTxHash, EventDeleteGroup))
	return map[string][]interface{}{
		k:  v,
		ek: ev,
	}
}

func (m *Module) handleLeaveGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, leaveGroup *storagetypes.EventLeaveGroup) []database.SQLStatement {
	g := &models.Group{
		OwnerAddress:    leaveGroup.Owner,
		GroupID:         leaveGroup.GroupId.BigInt().String(),
		GroupName:       leaveGroup.GroupName,
		AccountAddress:  leaveGroup.MemberAddress,
		UpdateAt:        block.Block.Height,
		UpdateTime:      block.Block.Time,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		Removed:         true,
	}
	k1, v1 := m.db.UpdateGroupMemberRemovedToSQL(ctx, g)
	k2, v2 := m.db.UpdateGroupMetaToSQL(ctx, g.GroupID, block.Block.Height, block.Block.Time, txHash, evmTxHash)
	ek, ev := m.db.SaveGroupEventToSQL(ctx, models.NewGroupEvent(g.GroupID, block.Block.Height, txHash, evmTxHash, EventLeaveGroup))
	return []database.SQLStatement{
		{SQL: k1, Vars: v1},
		{SQL: k2, Vars: v2},
		{SQL: ek, Vars: ev},
	}
}

func (m *Module) handleUpdateGroupMember(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, updateGroupMember *storagetypes.EventUpdateGroupMember) map[string][]interface{} {
	membersToAdd := updateGroupMember.MembersToAdd
	membersToDelete := updateGroupMember.MembersToDelete

	var membersToAddList []*models.Group
	res := make(map[string][]interface{})

	if len(membersToAdd) > 0 {
		for _, memberToAdd := range membersToAdd {
			groupItem := &models.Group{
				OwnerAddress:    updateGroupMember.Owner,
				GroupID:         updateGroupMember.GroupId.BigInt().String(),
				GroupName:       updateGroupMember.GroupName,
				AccountAddress:  memberToAdd.Member,
				OperatorAddress: updateGroupMember.Operator,
				CreateAt:        block.Block.Height,
				CreateTime:      block.Block.Time,
				CreateTxHash:    txHash,
				CreateEVMTxHash: evmTxHash,
				UpdateAt:        block.Block.Height,
				UpdateTime:      block.Block.Time,
				UpdateTxHash:    txHash,
				UpdateEVMTxHash: evmTxHash,
				Removed:         false,
				// ExpirationTime:  0,
			}
			if memberToAdd.ExpirationTime != nil {
				groupItem.ExpirationTime = *memberToAdd.ExpirationTime
			}
			membersToAddList = append(membersToAddList, groupItem)
		}
		k, v := m.db.CreateGroupToSQL(ctx, membersToAddList)
		res[k] = v
	}

	if len(membersToDelete) > 0 {
		groupItem := &models.Group{
			OperatorAddress: updateGroupMember.Operator,
			UpdateAt:        block.Block.Height,
			UpdateTime:      block.Block.Time,
			UpdateTxHash:    txHash,
			UpdateEVMTxHash: evmTxHash,
			Removed:         true,
		}
		accountAddresses := make([]string, 0, len(membersToDelete))
		for _, memberToDelete := range membersToDelete {
			accountAddresses = append(accountAddresses, memberToDelete)
		}
		k, v := m.db.BatchDeleteGroupMemberToSQL(ctx, groupItem, updateGroupMember.GroupId.BigInt().String(), accountAddresses)
		res[k] = v
	}

	// update group item
	groupItem := &models.Group{
		GroupID:         updateGroupMember.GroupId.BigInt().String(),
		UpdateAt:        block.Block.Height,
		UpdateTime:      block.Block.Time,
		UpdateTxHash:    txHash,
		UpdateEVMTxHash: evmTxHash,
		Removed:         false,
	}
	k, v := m.db.UpdateGroupToSQL(ctx, groupItem)
	res[k] = v
	ek, ev := m.db.SaveGroupEventToSQL(ctx, models.NewGroupEvent(groupItem.GroupID, block.Block.Height, txHash, evmTxHash, EventUpdateGroupMember))
	res[ek] = ev
	return res
}

func (m *Module) handleRenewGroupMember(ctx context.Context, block *tmctypes.ResultBlock, txHash, evmTxHash string, renewGroupMember *storagetypes.EventRenewGroupMember) []database.SQLStatement {
	stmts := make([]database.SQLStatement, 0, len(renewGroupMember.Members)+1)
	const renewSQL = `UPDATE "groups" SET expiration_time = ?, update_at = ?, update_time = ? WHERE account_address = ? AND group_id = ?`
	for _, e := range renewGroupMember.Members {
		var expirationTime interface{}
		if e.ExpirationTime != nil {
			expirationTime = e.ExpirationTime.UTC()
		}
		v := []interface{}{expirationTime, block.Block.Height, block.Block.Time, e.Member, renewGroupMember.GroupId.BigInt().String()}
		stmts = append(stmts, database.SQLStatement{SQL: renewSQL, Vars: v})
	}
	e := &models.GroupEvent{
		GroupID:   renewGroupMember.GroupId.BigInt().String(),
		Height:    block.Block.Height,
		TxHash:    txHash,
		EVMTxHash: evmTxHash,
		Event:     EventRenewGroupMember,
	}
	ek, ev := m.db.SaveGroupEventToSQL(ctx, e)
	stmts = append(stmts, database.SQLStatement{SQL: ek, Vars: ev})
	return stmts
}
