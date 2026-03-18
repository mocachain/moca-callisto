package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// Use types from moca repository
	paymenttypes "github.com/evmos/evmos/v12/x/payment/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	paymenttypes.RegisterCodec(cdc)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	paymenttypes.RegisterInterfaces(registry)
}

var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
