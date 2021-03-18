package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type rewardPoolSA struct {
	UserAddress     string `json:"userAddress"`
	TransactionType string `json:"transactionType"`
	Amount          string `json:"amount"`
	BlockTimestamp  int64  `json:"blockTimestamp"`
	BlockNumber     int64  `json:"blockNumber"`
	TxHash          string `json:"transactionHash"`
}

func (a *API) handleRewardPoolsStakingActions(c *gin.Context) {
	poolAddress, err := getQueryAddress(c, "poolAddress")
	if err != nil {
		Error(c, err)
		return
	}

	userAddress := strings.ToLower(c.DefaultQuery("userAddress", "all"))

	transactionType := strings.ToUpper(c.DefaultQuery("transactionType", "all"))

	if !checkRewardPoolTxType(transactionType) && transactionType != "ALL" {
		Error(c, errors.New("transaction type does not exist"))
		return
	}

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

	var parameters []interface{}
	parameters = append(parameters, poolAddress, limit, offset)

	query := `select 
    						user_address,
						    amount,
						    action_type,
						    block_timestamp,
						    included_in_block,
						    tx_hash
						from smart_yield_rewards_staking_actions 
						where pool_address = $1 %s %s 
						order by included_in_block desc ,
						         tx_index desc, 
						         log_index desc
						limit $2 offset  $3 `

	var userFilter, txTypeFilter string

	if userAddress != "all" {
		parameters = append(parameters, userAddress)
		userFilter = fmt.Sprintf("and user_address = $%d", len(parameters))
	}

	if transactionType != "ALL" {
		parameters = append(parameters, transactionType)
		txTypeFilter = fmt.Sprintf("and action_type = $%d", len(parameters))
	}

	query = fmt.Sprintf(query, userFilter, txTypeFilter)

	rows, err := a.db.Query(query, parameters...)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var transactions []rewardPoolSA
	for rows.Next() {
		var t rewardPoolSA
		err := rows.Scan(&t.UserAddress, &t.Amount, &t.TransactionType, &t.BlockTimestamp, &t.BlockNumber, &t.TxHash)
		if err != nil {
			Error(c, err)
			return
		}

		transactions = append(transactions, t)
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}
	count := len(transactions)

	OK(c, transactions, map[string]interface{}{
		"block": block,
		"count": count,
	})
}

func checkRewardPoolTxType(action string) bool {
	txType := [2]string{"JUNIOR_STAKE", "JUNIOR_UNSTAKE"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}
