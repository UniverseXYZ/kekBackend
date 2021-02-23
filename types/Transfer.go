package types

import (
	"math/big"
)

type Transfer struct {
	SYAddress        string
	ProtocolId       string
	TokenAddress     string
	From             string
	To               string
	Value            *big.Int
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}
