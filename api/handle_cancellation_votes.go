package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) CancellationVotesHandler(c *gin.Context) {
	proposalID := c.Param("proposalID")
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	supportFilter := strings.ToLower(c.DefaultQuery("support", ""))

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	var cancellationVotesList []types.Vote
	var rows *sql.Rows
	if supportFilter == "" {
		rows, err = a.db.Query(`select * from cancellation_proposal_votes($1) order by power desc offset $2 limit $3`, proposalID, offset, limit)
	} else {
		if supportFilter != "true" && supportFilter != "false" {
			BadRequest(c, errors.New("wrong value for support parameter"))
			return
		}
		rows, err = a.db.Query(`select * from cancellation_proposal_votes($1) where support = $4 order by power desc offset $2 limit $3`, proposalID, offset, limit, supportFilter)
	}

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			user           string
			support        bool
			blockTimestamp int64
			power          string
		)
		err := rows.Scan(&user, &support, &blockTimestamp, &power)
		if err != nil {
			Error(c, err)
			return
		}

		cancellationVote := types.Vote{
			User:           user,
			BlockTimestamp: blockTimestamp,
			Support:        support,
			Power:          power,
		}

		cancellationVotesList = append(cancellationVotesList, cancellationVote)
	}

	if len(cancellationVotesList) == 0 {
		NotFound(c)
		return
	}

	var count int
	if supportFilter == "" {
		err = a.db.QueryRow(`select count(*) from cancellation_proposal_votes($1)`, proposalID).Scan(&count)
	} else {
		err = a.db.QueryRow(`select count(*) from cancellation_proposal_votes($1) where support = $2`, proposalID, supportFilter).Scan(&count)
	}

	OK(c, cancellationVotesList, map[string]interface{}{"count": count})

}
