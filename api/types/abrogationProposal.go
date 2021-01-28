package types

type AbrogationProposal struct {
	ProposalID  uint64 `json:"proposalId"`
	Creator     string `json:"caller"`
	CreateTime  uint64 `json:"createTime"`
	Description string `json:"description"`
}
