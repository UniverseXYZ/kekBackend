package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleNotifications(c *gin.Context) {
	limit, err := getQueryLimit(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	page, err := getQueryPage(c)
	if err != nil {
		BadRequest(c, err)
	}

	offset := (page - 1) * limit

	query := `
				select 
					   target,
					   "type",
					   starts_on,
					   expires_on,
					   "message",
					   "metadata"
				from notifications
				where (target = 'system' %s) %s
				order by starts_on desc
				offset $1 limit $2`

	var parameters = []interface{}{offset, limit}

	target := strings.ToLower(c.DefaultQuery("target", "system"))
	var targetFilter string

	if target != "system" {
		parameters = append(parameters, target)
		targetFilter = fmt.Sprintf("or target = $%d", len(parameters))
	}

	timestamp, err := getQueryTimestamp(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	var timestampFilter string

	if timestamp > 0 {
		parameters = append(parameters, timestamp)
		timestampFilter = fmt.Sprintf("and starts_on < $%d and expires_on > $%d ", len(parameters), len(parameters))
	}

	query = fmt.Sprintf(query, targetFilter, timestampFilter)

	rows, err := a.db.Query(query, parameters...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var notifications []types.Notification
	for rows.Next() {
		var n types.Notification
		err := rows.Scan(&n.Target, &n.NotificationType, &n.StartsOn, &n.ExpiresOn, &n.Message, &n.Metadata)
		if err != nil {
			Error(c, err)
			return
		}
		notifications = append(notifications, n)
	}

	OK(c, notifications)

}
