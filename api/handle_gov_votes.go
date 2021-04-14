package api

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/api/types"
	"github.com/kekDAO/kekBackend/utils"
)

func (a *API) VotesHandler(c *gin.Context) {
	proposalIDString := utils.CleanUpHex(c.Param("proposalID"))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	supportFilter := strings.ToLower(c.DefaultQuery("support", ""))

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	proposalID, err := strconv.Atoi(proposalIDString)
	if err != nil {
		Error(c, err)
		return
	}

	var votesList []types.Vote
	var rows *sql.Rows

	if supportFilter == "" {
		rows, err = a.db.Query(`select * from proposal_votes($1)  order by power desc offset $2 limit $3`, proposalID, offset, limit)
	} else {
		if supportFilter != "true" && supportFilter != "false" {
			BadRequest(c, errors.New("wrong value for support parameter"))
			return
		}
		rows, err = a.db.Query(`select * from proposal_votes($1) where support = $4  order by power desc offset $2 limit $3`, proposalID, offset, limit, supportFilter)
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

		vote := types.Vote{
			User:           user,
			BlockTimestamp: blockTimestamp,
			Support:        support,
			Power:          power,
		}
		votesList = append(votesList, vote)
	}

	var count int
	if supportFilter == "" {
		err = a.db.QueryRow(`select count(*) from proposal_votes($1)`, proposalID).Scan(&count)
	} else {
		err = a.db.QueryRow(`select count(*) from proposal_votes($1) where support = $2`, proposalID, supportFilter).Scan(&count)
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, votesList, map[string]interface{}{"count": count, "block": block})
}
