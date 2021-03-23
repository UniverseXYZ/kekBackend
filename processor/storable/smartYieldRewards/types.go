package smartYieldRewards

import (
	"math/big"

	"github.com/barnbridge/barnbridge-backend/types"
)

type StakingAction struct {
	*types.Event
	UserAddress  string
	Amount       *big.Int
	BalanceAfter *big.Int
	ActionType   ActionType

	PoolAddress string
}

type ClaimEvent struct {
	*types.Event
	UserAddress string
	Amount      *big.Int

	PoolAddress string
}
