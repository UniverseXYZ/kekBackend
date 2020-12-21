package types

import (
	types2 "github.com/barnbridge/barnbridge-backend/types"
)

type Proposal struct {
	Id           uint64 `json:"proposal_ID,omitempty"`
	Proposer     string `json:"proposer,omitempty"`
	Description  string `json:"description,omitempty"`
	Title        string `json:"title,omitempty"`
	CreateTime   uint64 `json:"create_time,omitempty"`
	StartTime    uint64 `json:"start_time,omitempty"`
	Quorum       string `json:"quorum,omitempty"`
	Eta          uint64 `json:"eta,omitempty"`
	ForVotes     string `json:"for_votes,omitempty"`
	AgainstVotes string `json:"against_votes,omitempty"`
	Canceled     bool   `json:"canceled,omitempty"`
	Executed     bool   `json:"executed,omitempty"`

	Targets    types2.JSONStringArray `json:"targets,omitempty"`
	Values     types2.JSONStringArray `json:"values,omitempty"`
	Signatures types2.JSONStringArray `json:"signatures,omitempty"`
	Calldatas  types2.JSONStringArray `json:"calldatas,omitempty"`

	Timestamp int64 `json:"timestamp"`
}
