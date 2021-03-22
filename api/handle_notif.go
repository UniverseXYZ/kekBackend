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
		SELECT
		    "id",
			"target",
			"type",
			"starts_on",
			"expires_on",
			"message",
			"metadata"
		FROM
			"notifications"
		WHERE
			  $3 < "starts_on"
		  AND "starts_on" < EXTRACT(EPOCH FROM NOW())::bigint
		  AND EXTRACT(EPOCH FROM NOW())::bigint < "expires_on"
		  AND (
					  "target" = 'system' %s
				  )
		ORDER BY
			"starts_on"
		OFFSET $1 LIMIT $2
		;
	`

	var parameters = []interface{}{offset, limit}

	timestamp, err := getQueryTimestamp(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	parameters = append(parameters, timestamp)

	target := strings.ToLower(c.DefaultQuery("target", ""))
	var targetFilter string

	if target == "" {
		parameters = append(parameters, target)
		targetFilter = fmt.Sprintf(`OR target = $%d`, len(parameters))
	}

	query = fmt.Sprintf(query, targetFilter)

	rows, err := a.db.Query(query, parameters...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	defer rows.Close()

	var notifications []types.Notification
	for rows.Next() {
		var n types.Notification
		err := rows.Scan(&n.Id, &n.Target, &n.NotificationType, &n.StartsOn, &n.ExpiresOn, &n.Message, &n.Metadata)
		if err != nil {
			Error(c, err)
			return
		}
		notifications = append(notifications, n)
	}

	OK(c, notifications)
}
