//nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Expose the hooks if present
func (k Keeper) PostValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostValidatorCreated(ctx, valAddr)
	}
}

func (k Keeper) PostValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostValidatorModified(ctx, valAddr)
	}
}

func (k Keeper) PostValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostValidatorRemoved(ctx, consAddr, valAddr)
	}
}

func (k Keeper) PostValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostValidatorBonded(ctx, consAddr, valAddr)
	}
}

func (k Keeper) PostValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostValidatorPowerDidChange(ctx, consAddr, valAddr)
	}
}

func (k Keeper) PostValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
}

func (k Keeper) PostDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostDelegationCreated(ctx, delAddr, valAddr)
	}
}

func (k Keeper) PostDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostDelegationSharesModified(ctx, delAddr, valAddr)
	}
}

func (k Keeper) PostDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.PostDelegationRemoved(ctx, delAddr, valAddr)
	}
}
