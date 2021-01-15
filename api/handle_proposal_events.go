package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleProposalEvents(c *gin.Context) {
	proposalID := c.Param("proposalID")

	rows, err := a.db.Query(`select proposal_id ,caller,event_type,event_data,block_timestamp from governance_events where proposal_id = $1`, proposalID)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var eventsList []types.Event

	for rows.Next() {
		var event types.Event
		err := rows.Scan(&event.ProposalID, &event.Caller, &event.EventType, &event.Eta, &event.CreateTime)
		if err != nil {
			Error(c, err)
			return
		}

		eventsList = append(eventsList, event)
	}

	if len(eventsList) == 0 {
		NotFound(c)
		return
	}
	OK(c, eventsList)
}
