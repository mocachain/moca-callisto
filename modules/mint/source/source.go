package source

import (
	sdkmath "cosmossdk.io/math"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type Source interface {
	GetInflation(height int64) (sdkmath.LegacyDec, error)
	Params(height int64) (minttypes.Params, error)
}
