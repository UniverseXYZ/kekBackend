package api

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleStakingActions(c *gin.Context) {
	actionType := strings.ToUpper(c.DefaultQuery("type", "all"))
	user := strings.ToLower(c.DefaultQuery("user", ""))
	token := strings.ToLower(c.DefaultQuery("token", ""))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	tsString := c.DefaultQuery("timestamp", "0")
	direction := strings.ToLower(c.DefaultQuery("direction", "desc"))

	var timestamp int64

	if tsString != "0" {
		var err error
		timestamp, err = strconv.ParseInt(tsString, 10, 64)
		if err != nil {
			BadRequest(c, errors.New("unknown timestamp"))
			return
		}
	}

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
	%s %s %s %s
	order by block_timestamp %s
	offset $1
	limit $2`

	var actionFilter string
	var parameters = []interface{}{offset, limit}

	if actionType != "ALL" {
		parameters = append(parameters, actionType)
		actionFilter = fmt.Sprintf("and action_type = $%d", len(parameters))
	}

	var userFilter string
	if user != "" {
		parameters = append(parameters, user)
		userFilter = fmt.Sprintf("and user_address = $%d", len(parameters))
	}

	var tokenFilter string
	if token != "" {
		parameters = append(parameters, token)
		tokenFilter = fmt.Sprintf("and token_address = $%d", len(parameters))
	}

	var timestampFilter string

	if timestamp != 0 {
		if direction == "desc" {
			timestampFilter = fmt.Sprintf("and block_timestamp < %d", timestamp)
		} else {
			timestampFilter = fmt.Sprintf("and block_timestamp > %d", timestamp)
		}
	}

	query = fmt.Sprintf(query, actionFilter, userFilter, tokenFilter, timestampFilter, direction)
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
