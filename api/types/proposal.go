package types

import (
	types2 "github.com/barnbridge/barnbridge-backend/types"
)

type Proposal struct {
	Id          uint64 `json:"proposalId"`
	Proposer    string `json:"proposer"`
	Description string `json:"description"`
	Title       string `json:"title"`
	CreateTime  int64  `json:"createTime"`

	Targets    types2.JSONStringArray `json:"targets"`
	Values     types2.JSONStringArray `json:"values"`
	Signatures types2.JSONStringArray `json:"signatures"`
	Calldatas  types2.JSONStringArray `json:"calldatas"`

	BlockTimestamp      int64 `json:"blockTimestamp"`
	WarmUpDuration      int64 `json:"warmUpDuration"`
	ActiveDuration      int64 `json:"activeDuration"`
	QueueDuration       int64 `json:"queueDuration"`
	GracePeriodDuration int64 `json:"gracePeriodDuration"`
	AcceptanceThreshold int64 `json:"acceptanceThreshold"`
	MinQuorum           int64 `json:"minQuorum"`

	State         string `json:"state"`
	StateTimeLeft *int64 `json:"stateTimeLeft"`
	ForVotes      string `json:"forVotes"`
	AgainstVotes  string `json:"againstVotes"`
}

type ProposalLite struct {
	Id            uint64 `json:"proposalId"`
	Proposer      string `json:"proposer"`
	Description   string `json:"description"`
	Title         string `json:"title"`
	CreateTime    int64  `json:"createTime"`
	State         string `json:"state"`
	StateTimeLeft *int64 `json:"stateTimeLeft"`
	ForVotes      string `json:"forVotes"`
	AgainstVotes  string `json:"againstVotes"`
}
