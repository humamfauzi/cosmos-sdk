package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHookOnValidatorBonded(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, DefaultParams())
	addr := sdk.ConsAddress(addrs[0])
	keeper.postValidatorBonded(ctx, addr, nil)
	period := keeper.getValidatorSlashingPeriodForHeight(ctx, addr, ctx.BlockHeight())
	require.Equal(t, ValidatorSlashingPeriod{addr, ctx.BlockHeight(), 0, sdk.ZeroDec()}, period)
}

func TestHookOnValidatorBeginUnbonding(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, DefaultParams())
	addr := sdk.ConsAddress(addrs[0])
	keeper.postValidatorBonded(ctx, addr, nil)
	keeper.postValidatorBeginUnbonding(ctx, addr, addrs[0])
	period := keeper.getValidatorSlashingPeriodForHeight(ctx, addr, ctx.BlockHeight())
	require.Equal(t, ValidatorSlashingPeriod{addr, ctx.BlockHeight(), ctx.BlockHeight(), sdk.ZeroDec()}, period)
}
