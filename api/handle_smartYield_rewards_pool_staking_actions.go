package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type rewardPoolSA struct {
	UserAddress     string          `json:"userAddress"`
	TransactionType string          `json:"transactionType"`
	Amount          decimal.Decimal `json:"amount"`
	BlockTimestamp  int64           `json:"blockTimestamp"`
	BlockNumber     int64           `json:"blockNumber"`
	TxHash          string          `json:"transactionHash"`
}

func (a *API) handleRewardPoolsStakingActions(c *gin.Context) {
	poolAddress, err := getQueryAddress(c, "poolAddress")
	if err != nil {
		Error(c, err)
		return
	}

	rewardPool := state.RewardPoolByAddress(poolAddress)
	if rewardPool == nil {
		BadRequest(c, errors.New("invalid reward pool address"))
		return
	}

	smartYield := state.PoolBySmartYieldAddress(rewardPool.PoolTokenAddress)
	if smartYield == nil {
		Error(c, errors.New("could not find smart yield pool"))
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

	filters := new(Filters)
	filters.Add("pool_address", poolAddress)

	userAddress := c.DefaultQuery("userAddress", "all")
	if userAddress != "all" {
		filters.Add("user_address", utils.NormalizeAddress(userAddress))
	}

	transactionType := strings.ToUpper(c.DefaultQuery("transactionType", "all"))
	if transactionType != "ALL" {
		if !checkRewardPoolTxType(transactionType) {
			Error(c, errors.New("transaction type does not exist"))
			return
		}

		filters.Add("action_type", transactionType)
	}

	query, params := buildQueryWithFilter(`
		select
			user_address,
			amount,
			action_type,
			block_timestamp,
			included_in_block,
			tx_hash
		from smart_yield_rewards_staking_actions
		where %s
		order by included_in_block desc ,
				 tx_index desc,
				 log_index desc
		%s %s
	`,
		filters,
		&limit,
		&offset,
	)

	rows, err := a.db.Query(query, params...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(smartYield.UnderlyingDecimals))

	var transactions []rewardPoolSA
	for rows.Next() {
		var t rewardPoolSA
		err := rows.Scan(&t.UserAddress, &t.Amount, &t.TransactionType, &t.BlockTimestamp, &t.BlockNumber, &t.TxHash)
		if err != nil {
			Error(c, err)
			return
		}

		t.Amount = t.Amount.DivRound(tenPowDec, int32(smartYield.UnderlyingDecimals))

		transactions = append(transactions, t)
	}

	query, params = buildQueryWithFilter(`
		select
			count(*)
		from smart_yield_rewards_staking_actions
		where %s
		%s %s
	`,
		filters,
		nil,
		nil,
	)

	var count int64
	err = a.db.QueryRow(query, params...).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

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
