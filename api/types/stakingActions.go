package types

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
