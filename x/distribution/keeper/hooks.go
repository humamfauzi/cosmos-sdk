package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Create a new validator distribution record
func (k Keeper) postValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {

	// defensive check for existence
	if k.HasValidatorDistInfo(ctx, valAddr) {
		panic("validator dist info already exists (not cleaned up properly)")
	}

	height := ctx.BlockHeight()
	vdi := types.ValidatorDistInfo{
		OperatorAddr:            valAddr,
		FeePoolWithdrawalHeight: height,
		DelAccum:                types.NewTotalAccum(height),
		DelPool:                 types.DecCoins{},
		ValCommission:           types.DecCoins{},
	}
	k.SetValidatorDistInfo(ctx, vdi)
}

// Withdraw all validator rewards
func (k Keeper) postValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	// Move the validator's rewards from the global pool to the validator's pools
	// (dist info), but without actually withdrawing the rewards. This does not
	// need to happen during the genesis block.
	if ctx.BlockHeight() > 0 {
		if err := k.takeValidatorFeePoolRewards(ctx, valAddr); err != nil {
			panic(err)
		}
	}
}

// Withdraw all validator rewards
func (k Keeper) postValidatorBonded(ctx sdk.Context, valAddr sdk.ValAddress) {
	lastPower := k.stakeKeeper.GetLastValidatorPower(ctx, valAddr)
	if !lastPower.Equal(sdk.ZeroInt()) {
		panic("expected last power to be 0 for validator entering bonded state")
	}
	k.postValidatorModified(ctx, valAddr)
}

// Sanity check, very useful!
func (k Keeper) postValidatorPowerDidChange(ctx sdk.Context, valAddr sdk.ValAddress) {
	vi := k.GetValidatorDistInfo(ctx, valAddr)
	if vi.FeePoolWithdrawalHeight != ctx.BlockHeight() {
		panic(fmt.Sprintf("expected validator (%v) dist info FeePoolWithdrawalHeight to be updated to %v, but was %v.",
			valAddr.String(), ctx.BlockHeight(), vi.FeePoolWithdrawalHeight))
	}
}

// Withdrawal all validator distribution rewards and cleanup the distribution record
func (k Keeper) postValidatorRemoved(ctx sdk.Context, valAddr sdk.ValAddress) {
	k.RemoveValidatorDistInfo(ctx, valAddr)
}

//_________________________________________________________________________________________

// Create a new delegator distribution record
func (k Keeper) postDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) {

	ddi := types.DelegationDistInfo{
		DelegatorAddr:           delAddr,
		ValOperatorAddr:         valAddr,
		DelPoolWithdrawalHeight: ctx.BlockHeight(),
	}
	k.SetDelegationDistInfo(ctx, ddi)
}

// Withdrawal all validator rewards
func (k Keeper) postDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) {

	if err := k.WithdrawDelegationReward(ctx, delAddr, valAddr); err != nil {
		panic(err)
	}
}

// Withdrawal all validator distribution rewards and cleanup the distribution record
func (k Keeper) postDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) {

	k.RemoveDelegationDistInfo(ctx, delAddr, valAddr)
}

//_________________________________________________________________________________________

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ sdk.StakingHooks = Hooks{}

// New Validator Hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// nolint
func (h Hooks) PostValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.k.postValidatorCreated(ctx, valAddr)
}
func (h Hooks) PostValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.k.postValidatorModified(ctx, valAddr)
}
func (h Hooks) PostValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorRemoved(ctx, valAddr)
}
func (h Hooks) PostDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorModified(ctx, valAddr)
	h.k.postDelegationCreated(ctx, delAddr, valAddr)
}
func (h Hooks) PostDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorModified(ctx, valAddr)
	h.k.postDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h Hooks) PostDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.postDelegationRemoved(ctx, delAddr, valAddr)
}
func (h Hooks) PostValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorModified(ctx, valAddr)
}
func (h Hooks) PostValidatorBonded(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorBonded(ctx, valAddr)
}
func (h Hooks) PostValidatorPowerDidChange(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorPowerDidChange(ctx, valAddr)
}
