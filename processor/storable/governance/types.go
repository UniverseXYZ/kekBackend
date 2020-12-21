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
	Caller     *string
	Eta        *big.Int
	EventType  int
}

type Vote struct {
	BaseLog

	ProposalID *big.Int
	User       string
	Support    *bool
	Power      *big.Int
	Timestamp  int64
}

type VoteCanceled struct {
	BaseLog

	ProposalID *big.Int
	User       string
	Timestamp  int64
}

type CancellationProposal struct {
	BaseLog

	ProposalID big.Int
	CreateTime int64
	Creator    string
}
