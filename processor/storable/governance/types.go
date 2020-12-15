package governance

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Proposal struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
	StartTime    *big.Int
	Quorum       *big.Int
	Eta          *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	Canceled     bool
	Executed     bool
}

type BaseLog struct {
	LoggedBy         string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type ProposalActions struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}

type ProposalEvent struct {
	BaseLog

	ProposalID *big.Int
	EventType  int
}

type Vote struct {
	BaseLog

	ProposalID *big.Int
	User       string
	Support    *bool
	Canceled   bool
	Power      *big.Int
	Timestamp  int64
}
