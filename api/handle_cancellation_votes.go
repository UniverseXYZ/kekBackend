package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) CancellationVotesHandler(c *gin.Context) {
	proposalID := c.Param("proposalID")
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
	}

	var cancellationVotesList []types.Vote

	rows, err := a.core.DB().Query(`select * from cancellation_proposal_votes($1) order by power desc offset $2 limit $3`, proposalID, offset, limit)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
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
	err = a.core.DB().QueryRow(`select count(*) from cancellation_proposal_votes($1)`, proposalID).Scan(&count)

	OK(c, cancellationVotesList, map[string]interface{}{"count": count})

}
