package types

import (
	"math/big"
)

type Transfer struct {
	From             string
	To               string
	Value            *big.Int
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}
