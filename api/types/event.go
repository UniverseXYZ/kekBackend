package types

import (
	"github.com/barnbridge/barnbridge-backend/types"
)

type Event struct {
	ProposalID uint64           `json:"proposalId"`
	Caller     string           `json:"caller"`
	Eta        types.JSONObject `json:"eventData"`
	EventType  string           `json:"eventType"`
	CreateTime int64            `json:"createTime"`
	TxHash     string           `json:"txHash"`
}
