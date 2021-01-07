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

	rows, err := a.core.DB().Query(`select distinct user_id,
                first_value(support) over (partition by user_id order by block_timestamp desc) as support,
                first_value(block_timestamp) over (partition by user_id order by block_timestamp desc) as block_timestamp,
                power
	from governance_cancellation_votes
	where proposal_id = $1
  	and ( select count(*)
        from governance_cancellation_votes_canceled
        where governance_cancellation_votes_canceled.proposal_id = governance_cancellation_votes.proposal_id
        and governance_cancellation_votes_canceled.user_id = governance_cancellation_votes.user_id
        and governance_cancellation_votes_canceled.block_timestamp > governance_cancellation_votes.block_timestamp ) = 0 order by power offset $2 limit $3`, proposalID, offset, limit)

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

	OK(c, cancellationVotesList)

}
