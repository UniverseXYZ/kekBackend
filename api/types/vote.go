package types

type Vote struct {
	LoggedBy         string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64

	ProposalID uint64
	User       string
	Support    bool
	Power      int64
	Timestamp  int64
}

type VoteCanceled struct {
	LoggedBy         string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64

	ProposalID uint64
	User       string
	Timestamp  int64
}
