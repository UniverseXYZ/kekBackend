package types

import (
	"github.com/barnbridge/barnbridge-backend/types"
)

type Notification struct {
	Target           string           `json:"target"`
	NotificationType string           `json:"notificationType"`
	TriggeredByBlock int64            `json:"triggeredByBlock"`
	StartsOn         int64            `json:"startsOn"`
	ExpiresOn        int64            `json:"expiresOn"`
	Message          string           `json:"message"`
	Metadata         types.JSONObject `json:"metadata"`
}
