package api

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) VotesHandler(c *gin.Context) {
	proposalIDString := utils.CleanUpHex(c.Param("proposalID"))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
	}

	proposalID, err := strconv.Atoi(proposalIDString)
	if err != nil {
		Error(c, err)
	}

	var votesList []types.Vote

	rows, err := a.core.DB().Query(`select distinct user_id,
                first_value(support) over (partition by user_id order by block_timestamp desc) as support,
                first_value(block_timestamp) over (partition by user_id order by block_timestamp desc) as block_timestamp,
                power
	from governance_votes
	where proposal_id = $1 
  	and ( select count(*)
        from governance_votes_canceled
        where governance_votes_canceled.proposal_id = governance_votes.proposal_id
        and governance_votes_canceled.user_id = governance_votes.user_id
        and governance_votes_canceled.block_timestamp > governance_votes.block_timestamp ) = 0 order by power desc offset $2 limit $3`, proposalID, offset, limit)

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

		vote := types.Vote{
			User:           user,
			BlockTimestamp: blockTimestamp,
			Support:        support,
			Power:          power,
		}
		votesList = append(votesList, vote)
	}
	if len(votesList) == 0 {
		NotFound(c)
		return
	}

	OK(c, votesList)
}
