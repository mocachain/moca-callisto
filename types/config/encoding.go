package config

import (
	"cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/types/module"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	ethermint "github.com/evmos/evmos/v12/encoding"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
)

// MakeEncodingConfig creates an EncodingConfig to properly handle all the messages
func MakeEncodingConfig(managers []module.BasicManager) func() params.EncodingConfig {
	return func() params.EncodingConfig {
		encodingConfig := ethermint.MakeConfig()
		evmtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		govv1.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		govv1beta1.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		return params.EncodingConfig{
			Amino:             encodingConfig.Amino,
			Codec:             encodingConfig.Codec,
			InterfaceRegistry: encodingConfig.InterfaceRegistry,
			TxConfig:          encodingConfig.TxConfig,
		}
	}
}

// mergeBasicManagers merges the given managers into a single module.BasicManager
func mergeBasicManagers(managers []module.BasicManager) module.BasicManager {
	var union = module.BasicManager{}
	for _, manager := range managers {
		for k, v := range manager {
			union[k] = v
		}
	}
	return union
}
