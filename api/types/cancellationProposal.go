package types

type CancellationProposal struct {
	ProposalID uint64 `json:"proposal_id"`
	Creator    string `json:"caller"`
	CreateTime uint64 `json:"create_time"`
}
