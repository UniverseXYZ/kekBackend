package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/kekDAO/kekBackend/api/types"
)

func (a *API) handleProposalEvents(c *gin.Context) {
	proposalID := c.Param("proposalID")

	rows, err := a.db.Query(`select proposal_id ,caller,event_type,event_data,block_timestamp, tx_hash from governance_events where proposal_id = $1`, proposalID)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var eventsList []types.Event

	for rows.Next() {
		var event types.Event
		err := rows.Scan(&event.ProposalID, &event.Caller, &event.EventType, &event.Eta, &event.CreateTime, &event.TxHash)
		if err != nil {
			Error(c, err)
			return
		}

		eventsList = append(eventsList, event)
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, eventsList, map[string]interface{}{
		"block": block,
	})
}
