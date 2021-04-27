package types

import (
	"github.com/kekDAO/kekBackend/types"
)

type Event struct {
	ProposalID uint64           `json:"proposalId"`
	Caller     string           `json:"caller"`
	Eta        types.JSONObject `json:"eventData"`
	EventType  string           `json:"eventType"`
	CreateTime int64            `json:"createTime"`
	TxHash     string           `json:"txHash"`
}
