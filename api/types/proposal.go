package types

import (
	types2 "github.com/barnbridge/barnbridge-backend/types"
)

type Proposal struct {
	Id          uint64 `json:"proposal_id,omitempty"`
	Proposer    string `json:"proposer,omitempty"`
	Description string `json:"description,omitempty"`
	Title       string `json:"title,omitempty"`
	CreateTime  uint64 `json:"create_time,omitempty"`

	Targets    types2.JSONStringArray `json:"targets,omitempty"`
	Values     types2.JSONStringArray `json:"values,omitempty"`
	Signatures types2.JSONStringArray `json:"signatures,omitempty"`
	Calldatas  types2.JSONStringArray `json:"calldatas,omitempty"`

	BlockTimestamp int64 `json:"block_timestamp"`
}
