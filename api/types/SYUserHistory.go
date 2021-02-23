package types

import (
	"github.com/shopspring/decimal"
)

type SYUserHistory struct {
	ProtocolId             string          `json:"protocolId"`
	Pool                   string          `json:"pool"`
	UnderlyingTokenAddress string          `json:"underlyingTokenAddress"`
	Amount                 decimal.Decimal `json:"amount"`
	Tranche                string          `json:"tranche"`
	TransactionType        string          `json:"transactionType"`
	TransactionHash        string          `json:"transactionHash"`
	BlockTimestamp         int64           `json:"blockTimestamp"`
	BlockNumber            int64           `json:"blockNumber"`
}
