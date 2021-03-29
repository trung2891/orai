package subscribe

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/tendermint/tendermint/libs/log"
)

// AIOracleConfig is the extra config required for wasm
type AIOracleConfig struct {
	BroadcastTimeout time.Duration
	RPCPollInterval  time.Duration
	MaxTry           uint64
	Txf              tx.Factory
	AllowLogLevel    log.Option
	Fees             sdk.Coins
	ErrExit          bool
}

func DefaultAIOracleConfig() *AIOracleConfig {
	return &AIOracleConfig{
		BroadcastTimeout: time.Minute * 5,
		RPCPollInterval:  time.Second,
		MaxTry:           5,
		AllowLogLevel:    log.AllowInfo(),
		ErrExit:          false,
		Fees:             sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(int64(5000)))),
	}
}