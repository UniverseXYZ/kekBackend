package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type StakingAction struct {
	UserAddress     string          `json:"userAddress"`
	TokenAddress    string          `json:"tokenAddress"`
	Amount          decimal.Decimal `json:"amount"`
	TransactionHash string          `json:"transactionHash"`
	ActionType      string          `json:"actionType"`
	BlockTimestamp  int64           `json:"blockTimestamp"`
}
type Pool struct {
	Tokens                []string
	EpochDelayFromStaking int64
}

type Chart map[time.Time]Aggregate

type Aggregate struct {
	SumDeposits    decimal.Decimal `json:"sumDeposits"`
	SumWithdrawals decimal.Decimal `json:"sumWithdrawals"`
}
