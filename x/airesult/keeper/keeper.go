package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/oraichain/orai/x/airesult/types"
)

// Keeper of the provider store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	supplyKeeper     types.SupplyKeeper
	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	distrKeeper      types.DistrKeeper
	ProviderKeeper   types.ProviderKeeper
	webSocketKeeper  types.WebSocketKeeper
	aiRequestKeeper  types.AIRequestKeeper
	feeCollectorName string
}

// NewKeeper creates a provider keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, aiResultSubSpace params.Subspace, supplyKeeper types.SupplyKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, distrKeeper types.DistrKeeper, providerKeeper types.ProviderKeeper, socketKeeper types.WebSocketKeeper, aiRequestKeeper types.AIRequestKeeper, feeCollectorName string) Keeper {
	if !aiResultSubSpace.HasKeyTable() {
		// register parameters of the provider module into the param space
		aiResultSubSpace = aiResultSubSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       aiResultSubSpace,
		supplyKeeper:     supplyKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		distrKeeper:      distrKeeper,
		ProviderKeeper:   providerKeeper,
		webSocketKeeper:  socketKeeper,
		aiRequestKeeper:  aiRequestKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//IsNamePresent checks if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// GetParam returns the parameter as specified by key as an uint64.
func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.paramSpace.Get(ctx, key, &res)
	return res
}

// SetParam saves the given key-value parameter to the store.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.paramSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}