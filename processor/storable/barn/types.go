package barn

import (
	"math/big"
)

type Withdraw struct {
	AmountWithdrew *big.Int
	AmountLeft     *big.Int
	User           string
}

type Deposit struct {
	Amount     *big.Int
	NewBalance *big.Int
	User       string
}

type BaseLog struct {
	LoggedBy         string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type Lock struct {
	BaseLog

	User        string
	LockedUntil *big.Int
}

type StakingAction struct {
	BaseLog

	UserAddress  string
	ActionType   int
	Amount       string
	BalanceAfter string
}

type DelegateAction struct {
	BaseLog

	Sender     string
	Receiver   string
	ActionType int
}

type DelegateChange struct {
	BaseLog

	ActionType          int
	Sender              string
	Receiver            string
	Amount              *big.Int
	ToNewDelegatedPower *big.Int
}
