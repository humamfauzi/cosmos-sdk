package slashing

import (
	"time"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) postValidatorBonded(ctx sdk.Context, address sdk.ConsAddress, _ sdk.ValAddress) {
	// Update the signing info start height or create a new signing info
	_, found := k.getValidatorSigningInfo(ctx, address)
	if !found {
		signingInfo := ValidatorSigningInfo{
			StartHeight:         ctx.BlockHeight(),
			IndexOffset:         0,
			JailedUntil:         time.Unix(0, 0),
			MissedBlocksCounter: 0,
		}
		k.SetValidatorSigningInfo(ctx, address, signingInfo)
	}

	// Create a new slashing period when a validator is bonded
	slashingPeriod := ValidatorSlashingPeriod{
		ValidatorAddr: address,
		StartHeight:   ctx.BlockHeight(),
		EndHeight:     0,
		SlashedSoFar:  sdk.ZeroDec(),
	}
	k.addOrUpdateValidatorSlashingPeriod(ctx, slashingPeriod)
}

// Mark the slashing period as having ended when a validator begins unbonding
func (k Keeper) postValidatorBeginUnbonding(ctx sdk.Context, address sdk.ConsAddress, _ sdk.ValAddress) {
	slashingPeriod := k.getValidatorSlashingPeriodForHeight(ctx, address, ctx.BlockHeight())
	slashingPeriod.EndHeight = ctx.BlockHeight()
	k.addOrUpdateValidatorSlashingPeriod(ctx, slashingPeriod)
}

// When a validator is created, add the address-pubkey relation.
func (k Keeper) postValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	validator := k.validatorSet.Validator(ctx, valAddr)
	k.addPubkey(ctx, validator.GetConsPubKey())
}

// When a validator is removed, delete the address-pubkey relation.
func (k Keeper) postValidatorRemoved(ctx sdk.Context, address sdk.ConsAddress) {
	k.deleteAddrPubkeyRelation(ctx, crypto.Address(address))
}

//_________________________________________________________________________________________

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ sdk.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) PostValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorBonded(ctx, consAddr, valAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) PostValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.postValidatorBeginUnbonding(ctx, consAddr, valAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) PostValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) {
	h.k.postValidatorRemoved(ctx, consAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) PostValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.k.postValidatorCreated(ctx, valAddr)
}

// nolint - unused hooks
func (h Hooks) PostValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) PostValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) PostDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) PostDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) PostDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
