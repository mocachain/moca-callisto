package types

import (
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/simapp/params"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmosapp "github.com/evmos/evmos/v12/app"
	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v4/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	vgtypes "github.com/evmos/evmos/v12/x/virtualgroup/types"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	permissionsource "github.com/forbole/bdjuno/v4/modules/permission/source"
	remotepermissionsource "github.com/forbole/bdjuno/v4/modules/permission/source/remote"
	permissiontypes "github.com/evmos/evmos/v12/x/permission/types"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	spsource "github.com/forbole/bdjuno/v4/modules/sp/source"
	remotespsource "github.com/forbole/bdjuno/v4/modules/sp/source/remote"
	sptypes "github.com/evmos/evmos/v12/x/sp/types"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
	storagesource "github.com/forbole/bdjuno/v4/modules/storage/source"
	remotestoragesource "github.com/forbole/bdjuno/v4/modules/storage/source/remote"
	storagetypes "github.com/evmos/evmos/v12/x/storage/types"
	vgsource "github.com/forbole/bdjuno/v4/modules/virtualgroup/source"
	remotevgsource "github.com/forbole/bdjuno/v4/modules/virtualgroup/source/remote"
	nodeconfig "github.com/forbole/juno/v5/node/config"
	"github.com/forbole/juno/v5/node/local"
	"github.com/forbole/juno/v5/node/remote"
)

type Sources struct {
	BankSource       banksource.Source
	DistrSource      distrsource.Source
	GovSource  govsource.Source
	MintSource mintsource.Source
	SlashingSource   slashingsource.Source
	StakingSource    stakingsource.Source
	StorageSource    storagesource.Source
	SpSource         spsource.Source
	VGSource         vgsource.Source
	PermissionSource permissionsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return nil, fmt.Errorf("local mode is currently not supported")

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := evmosapp.NewEvmos(
		log.NewNopLogger(), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, nil, nil, nil,
	)

	sources := &Sources{
		BankSource:      localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:     localdistrsource.NewSource(source, distrkeeper.Querier{Keeper: app.DistrKeeper}),
		GovSource:  localgovsource.NewSource(source, govkeeper.NewQueryServer(&app.GovKeeper), nil),
		MintSource: localmintsource.NewSource(source, minttypes.QueryServer(nil)), // no MintKeeper available in evmosApp
		SlashingSource:  localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:   localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(app, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(app, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:      remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:     remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:  remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn), govtypesv1beta1.NewQueryClient(source.GrpcConn)),
		MintSource: remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:  remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:   remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		StorageSource: remotestoragesource.NewSource(
			source,
			storagetypes.NewQueryClient(source.GrpcConn),
		),
		SpSource: remotespsource.NewSource(
			source,
			sptypes.NewQueryClient(source.GrpcConn),
		),
		VGSource: remotevgsource.NewSource(
			source,
			vgtypes.NewQueryClient(source.GrpcConn),
		),
		PermissionSource: remotepermissionsource.NewSource(
			source,
			permissiontypes.NewQueryClient(source.GrpcConn),
		),
	}, nil
}
