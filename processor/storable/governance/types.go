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
	ProposalParameters
}
type ProposalParameters struct {
	WarmUpDuration      *big.Int
	ActiveDuration      *big.Int
	QueueDuration       *big.Int
	GracePeriodDuration *big.Int
	AcceptanceThreshold *big.Int
	MinQuorum           *big.Int
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
	Caller     common.Address
	Eta        *big.Int
	EventType  ActionType
}

type Vote struct {
	BaseLog

	ProposalID *big.Int
	User       string
	Support    bool
	Power      *big.Int
	Timestamp  int64
}

type VoteCanceled struct {
	BaseLog

	ProposalID *big.Int
	User       string
	Timestamp  int64
}

type AbrogationProposal struct {
	BaseLog

	ProposalID  *big.Int
	CreateTime  int64
	Caller      common.Address
	Description string
}
