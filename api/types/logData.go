package types

type LogData struct {
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex int64  `json:"transactionIndex"`
	LogIndex         int64  `json:"logIndex"`
	BlockTimestamp   int64  `json:"blockTimestamp"`
	BlockNumber      int64  `json:"blockNumber"`
}
