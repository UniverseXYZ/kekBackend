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

	BlockTimestamp      int64  `json:"block_timestamp"`
	WarmUpDuration      uint64 `json:"warm_up_duration"`
	ActiveDuration      uint64 `json:"active_duration"`
	QueueDuration       uint64 `json:"queue_duration"`
	GracePeriodDuration uint64 `json:"grace_period_duration"`
	AcceptanceThreshold uint64 `json:"acceptance_threshold"`
	MinQuorum           uint64 `json:"min_quorum"`
	State               string `json:"state"`
}
