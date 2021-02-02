package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleStakingActions(c *gin.Context) {
	actionType := strings.ToUpper(c.DefaultQuery("type", "all"))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	if actionType != "ALL" && (actionType != "DEPOSIT" && actionType != "WITHDRAW") {
		BadRequest(c, errors.New("unknown state"))
		return
	}

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	query := `select 
				tx_hash,
				tx_index,
				log_index,
				user_address,
				token_address,
				amount ,
				action_type,
				block_timestamp,
				included_in_block
	from yield_farming_actions
	where 1=1 
	%s %s %s 
	order by block_timestamp desc
	offset $1
	limit $2`

	var actionFilter string
	var parameters = []interface{}{offset, limit}

	if actionType != "ALL" {
		parameters = append(parameters, actionType)
		actionFilter = fmt.Sprintf("and action_type = $%d", len(parameters))
	}

	query = fmt.Sprintf(query, actionFilter)
	spew.Dump(query)
	rows, err := a.db.Query(query, parameters...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var stakingActions []types.StakingAction
	for rows.Next() {
		var (
			user             string
			token            string
			amount           int64
			transactionHash  string
			transactionIndex int64
			logIndex         int64
			actionType       string
			blockTimestamp   int64
			blockNumber      int64
		)
		err := rows.Scan(&transactionHash, &transactionIndex, &logIndex, &user, &token, &amount, &actionType, &blockTimestamp, &blockNumber)
		if err != nil {
			Error(c, err)
			return
		}

		stakingAction := types.StakingAction{
			UserAddress:      user,
			TokenAddress:     token,
			Amount:           amount,
			TransactionHash:  transactionHash,
			TransactionIndex: transactionIndex,
			LogIndex:         logIndex,
			ActionType:       actionType,
			BlockTimestamp:   blockTimestamp,
			BlockNumber:      blockNumber,
		}
		stakingActions = append(stakingActions, stakingAction)

	}

	count := len(stakingActions)

	OK(c, stakingActions, map[string]interface{}{"count": count})

}
