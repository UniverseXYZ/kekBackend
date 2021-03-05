package smartYieldRewards

import (
	"math/big"

	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/types"
)

type StakingAction struct {
	*types.Event
	UserAddress  string
	Amount       decimal.Decimal
	BalanceAfter decimal.Decimal
	ActionType   ActionType

	PoolAddress string
}

type ClaimEvent struct {
	*types.Event
	UserAddress string
	Amount      *big.Int

	PoolAddress string
}
