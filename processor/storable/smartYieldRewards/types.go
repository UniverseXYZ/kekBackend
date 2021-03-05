package smartYieldRewards

import (
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/types"
)

type StakingEvent struct {
	*types.Event
	UserAddress  string
	Amount       decimal.Decimal
	BalanceAfter decimal.Decimal
	ActionType   StakingAction
}

type ClaimEvent struct {
	*types.Event
	UserAddress string
	Amount      decimal.Decimal
}
