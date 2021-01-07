package types

import (
	"github.com/barnbridge/barnbridge-backend/types"
)

type Event struct {
	ProposalID uint64           `json:"proposal_id"`
	Caller     string           `json:"caller"`
	Eta        types.JSONObject `json:"event_data"`
	EventType  string           `json:"event_type"`
	CreateTime uint64           `json:"create_time"`
}
