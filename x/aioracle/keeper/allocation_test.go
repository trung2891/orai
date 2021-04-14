package keeper_test

import (
	"testing"

	"github.com/oraichain/orai/x/aioracle"
	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	aioracletypes "github.com/oraichain/orai/x/aioracle/types"
)

var (
	PKS = simapp.CreateTestPubKeys(5)

	valConsPk1 = PKS[0]
	valConsPk2 = PKS[1]
	valConsPk3 = PKS[2]

	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())

	authAcc  = authtypes.NewEmptyModuleAccount("auth")
	distrAcc = authtypes.NewEmptyModuleAccount("distribution")
	bankAcc  = authtypes.NewEmptyModuleAccount("bank")
)

func TestAllocateTokensToManyValidators(t *testing.T) {

	// define static variables

	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrs(app, ctx, 10, sdk.NewInt(10000000000))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 10% commission and 300,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[7], valConsPk1, sdk.NewInt(30000000), true)

	// create second validator with 10% commission and 250,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[8], valConsPk2, sdk.NewInt(25000000), true)

	// create second validator with 10% commission and 150,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[9], valConsPk3, sdk.NewInt(15000000), true)

	abciValA := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   30,
	}
	abciValB := abci.Validator{
		Address: valConsPk2.Address(),
		Power:   25,
	}
	abciValC := abci.Validator{
		Address: valConsPk3.Address(),
		Power:   15,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[7]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[8]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[9]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[7]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[8]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[9]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[7]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[8]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[9]).Rewards.IsZero())

	// allocate tokens as if both had voted and second was proposer
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)))
	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	require.NotNil(t, feeCollector)

	err := app.BankKeeper.SetBalances(ctx, feeCollector.GetAddress(), fees)

	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.GetModuleAccount(ctx, auth.FeeCollectorName))

	votes := []abci.VoteInfo{
		{
			Validator:       abciValA,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValB,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValC,
			SignedLastBlock: true,
		},
	}

	testKeeper := aioraclekeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	id := ksuid.New().String()
	testKeeper.SetAIOracle(ctx, id, &aioracle.AIOracle{RequestID: id})

	aioracleTest, err := testKeeper.GetAIOracle(ctx, id)

	require.NoError(t, err)
	require.Equal(t, &aioracle.AIOracle{RequestID: id}, aioracleTest)

	// init reward
	reward := aioracletypes.DefaultReward(0)

	// init data sources
	firstDataSource := types.NewDataSourceResult(&types.EntryPoint{"", []string{}, addrs[1], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{}, "good")
	secondDataSource := types.NewDataSourceResult(&types.EntryPoint{"", []string{}, addrs[2], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{}, "good")
	thirdDataSource := types.NewDataSourceResult(&types.EntryPoint{"", []string{}, addrs[3], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{}, "good")

	// since there are three validators, we need to loop 3 times to add data sources
	for i := 0; i < 3; i++ {
		reward.DatasourceResults = append(reward.DatasourceResults, firstDataSource)
		reward.DatasourceResults = append(reward.DatasourceResults, secondDataSource)
		reward.DatasourceResults = append(reward.DatasourceResults, thirdDataSource)
	}

	reward.ProviderFees = reward.ProviderFees.Add(firstDataSource.GetEntryPoint().GetProviderFees()...).Add(secondDataSource.GetEntryPoint().GetProviderFees()...).Add(thirdDataSource.GetEntryPoint().GetProviderFees()...)

	// set validators
	rewardRatio := sdk.NewDec(int64(1)).Sub(sdk.NewDecWithPrec(int64(60), 2))
	valFees, _ := sdk.NewDecCoinsFromCoins(reward.ProviderFees...).MulDec(rewardRatio).TruncateDecimal()
	validatorA := &aioracletypes.Validator{valAddrs[7], abciValA.Power, "active"}
	validatorB := &aioracletypes.Validator{valAddrs[8], abciValB.Power, "active"}
	validatorC := &aioracletypes.Validator{valAddrs[9], abciValC.Power, "active"}
	// should be 15.6 ORAI
	reward.ValidatorFees = reward.ValidatorFees.Add(valFees...).Add(valFees...).Add(valFees...)
	reward.TotalPower = reward.TotalPower + validatorA.VotingPower + validatorB.VotingPower + validatorC.VotingPower
	reward.Validators = append(reward.Validators, *validatorA)
	reward.Validators = append(reward.Validators, *validatorB)
	reward.Validators = append(reward.Validators, *validatorC)

	tempReward := reward.ProviderFees
	// 39 ORAI
	reward.ProviderFees = reward.ProviderFees.Add(tempReward...).Add(tempReward...)

	// set reward
	testKeeper.SetReward(ctx, reward)

	t.Logf("balance of provider 1: %v\n", app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	t.Logf("balance of provider 2: %v\n", app.BankKeeper.GetBalance(ctx, addrs[1], sdk.DefaultBondDenom))

	t.Logf("balance of provider 3: %v\n", app.BankKeeper.GetBalance(ctx, addrs[2], sdk.DefaultBondDenom))

	t.Logf("balance of provider 4: %v\n", app.BankKeeper.GetBalance(ctx, addrs[3], sdk.DefaultBondDenom))

	t.Logf("balance of provider 5: %v\n", app.BankKeeper.GetBalance(ctx, addrs[4], sdk.DefaultBondDenom))

	testKeeper.AllocateTokens(ctx, votes, 1)
	rewardObj, err := testKeeper.GetReward(ctx, 1-1)

	require.NoError(t, err)

	// confirm that each validator fee is 5.2 ORAI => remove decimal to 5 ORAI
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(5)}}, valFees)

	// provider fees must equal 13 * 3 ORAI (3 validators)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDecWithPrec(39, 0)}}, sdk.NewDecCoinsFromCoins(rewardObj.ProviderFees...))

	// validator fees must equal 15.6 ORAI => 15 ORAI
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(15)}}, rewardObj.ValidatorFees)

	// reward must equal to 54.6 ORAI => 54 ORAI
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(54)}}, rewardObj.ProviderFees.Add(rewardObj.ValidatorFees...))

	// total power must equal
	require.Equal(t, int64(rewardObj.TotalPower), validatorA.VotingPower+validatorB.VotingPower+validatorC.VotingPower)
	require.Equal(t, int64(rewardObj.TotalPower), int64(70))

	t.Logf("after allocation\n")

	t.Logf("outstanding reward of validators: %v\n", app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[8]))

	t.Logf("balance of provider 1: %v\n", app.BankKeeper.GetBalance(ctx, addrs[1], sdk.DefaultBondDenom))

	t.Logf("balance of provider 2: %v\n", app.BankKeeper.GetBalance(ctx, addrs[2], sdk.DefaultBondDenom))

	t.Logf("balance of provider 3: %v\n", app.BankKeeper.GetBalance(ctx, addrs[3], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000015)), app.BankKeeper.GetBalance(ctx, addrs[1], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000015)), app.BankKeeper.GetBalance(ctx, addrs[2], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000009)), app.BankKeeper.GetBalance(ctx, addrs[3], sdk.DefaultBondDenom))

	// validate outstanding validator reward
	powerFractionVal7 := sdk.NewDec(30000000).QuoTruncate(sdk.NewDec(rewardObj.TotalPower))
	powerFractionVal8 := sdk.NewDec(25000000).QuoTruncate(sdk.NewDec(rewardObj.TotalPower))
	powerFractionVal9 := sdk.NewDec(15000000).QuoTruncate(sdk.NewDec(rewardObj.TotalPower))

	valRewardDec7, _ := sdk.NewDecCoinsFromCoins(rewardObj.ValidatorFees...).QuoDec(sdk.NewDec(int64(len(rewardObj.Validators)))).MulDec(powerFractionVal7).QuoDec(sdk.NewDec(1000000)).TruncateDecimal()
	valRewardDec8, _ := sdk.NewDecCoinsFromCoins(rewardObj.ValidatorFees...).QuoDec(sdk.NewDec(int64(len(rewardObj.Validators)))).MulDec(powerFractionVal8).QuoDec(sdk.NewDec(1000000)).TruncateDecimal()
	valRewardDec9, _ := sdk.NewDecCoinsFromCoins(rewardObj.ValidatorFees...).QuoDec(sdk.NewDec(int64(len(rewardObj.Validators)))).MulDec(powerFractionVal9).QuoDec(sdk.NewDec(1000000)).TruncateDecimal()

	realReward7, _ := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[7]).Rewards.TruncateDecimal()
	realReward8, _ := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[8]).Rewards.TruncateDecimal()
	realReward9, _ := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[9]).Rewards.TruncateDecimal()

	t.Logf("outstanding reward of validators: %v\n", valRewardDec7)

	// have to multiply with 1000000 because power fraction gets voting power raw (not * 10^-6)
	require.Equal(t, valRewardDec7, realReward7)
	require.Equal(t, valRewardDec8, realReward8)
	require.Equal(t, valRewardDec9, realReward9)
}
