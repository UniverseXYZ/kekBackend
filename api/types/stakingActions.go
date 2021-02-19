package types

import (
	"github.com/shopspring/decimal"
)

type StakingAction struct {
	UserAddress      string `json:"user"`
	TokenAddress     string `json:"token"`
	Amount           int64  `json:"amount"`
	TransactionHash  string `json:"txHash"`
	TransactionIndex int64  `json:"transactionIndex"`
	LogIndex         int64  `json:"logIndex"`
	ActionType       string `json:"type"`
	BlockTimestamp   int64  `json:"blockTimestamp"`
	BlockNumber      int64  `json:"blockNumber"`
}
type Pool struct {
	Tokens                []string
	EpochDelayFromStaking int64
}

type Chart map[string]Aggregate

type Aggregate struct {
	SumDeposits    decimal.Decimal `json:"sumDeposits"`
	SumWithdrawals decimal.Decimal `json:"sumWithdrawals"`
}
