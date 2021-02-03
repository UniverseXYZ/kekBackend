package api

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleStakingActionsList(c *gin.Context) {
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
	stakingActions, err := a.execQuery(query, parameters)
	if err != nil {
		Error(c, err)
		return
	}

	count := len(stakingActions)

	OK(c, stakingActions, map[string]interface{}{"count": count})

}

func (a *API) handleStakinsActionsChart(c *gin.Context) {
	period := c.DefaultQuery("period", "allActions")
	pool := strings.ToLower(c.DefaultQuery("pool", "stable"))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	actionType := strings.ToUpper(c.DefaultQuery("type", "ALL"))

	if actionType != "ALL" && (actionType != "DEPOSIT" && actionType != "WITHDRAW") {
		BadRequest(c, errors.New("unknown state"))
		return
	}

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	p, exists := types.Pools[pool]
	if !exists {
		Error(c, errors.New("could not find pool"))
		return
	}

	var epoch int64
	if period == "" || period == "allActions" {
		epoch = -1
	} else {
		epoch, err = strconv.ParseInt(period, 10, 64)
		if err != nil {
			Error(c, errors.New("invalid period"))
			return
		}
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
	%s 
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
	allActions, err := a.execQuery(query, parameters)
	if err != nil {
		Error(c, err)
		return
	}

	actions := filterActionsByPool(allActions, p)

	buckets := bucketActionsByEpoch(actions)
	var chart = make(types.Chart)
	if epoch == -1 {
		for epoch, actions := range buckets {
			key := strconv.FormatInt(epoch-p.EpochDelayFromStaking, 10)
			chart[key] = calcSum(actions)
		}

		count := len(chart)
		OK(c, chart, map[string]interface{}{"count": count})

		return
	}

	epoch = epoch + p.EpochDelayFromStaking
	actions = buckets[epoch]
	bucketsByDay := bucketActionsByDay(actions)

	for day, actions := range bucketsByDay {
		chart[day] = calcSum(actions)
	}

	count := len(chart)

	OK(c, chart, map[string]interface{}{"count": count})
	return
}

func (a *API) execQuery(query string, parameters []interface{}) ([]types.StakingAction, error) {
	rows, err := a.db.Query(query, parameters...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
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
			return nil, err
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
	return stakingActions, nil
}

func calcEpochByTimestamp(ts int64) int64 {
	if ts < types.Epoch1StartUnix {
		return 0
	}

	return (ts-types.Epoch1StartUnix)/types.EpochDuration + 1
}

func bucketActionsByEpoch(actions []types.StakingAction) map[int64][]types.StakingAction {
	res := make(map[int64][]types.StakingAction)

	for _, a := range actions {
		epoch := calcEpochByTimestamp(a.BlockTimestamp)
		res[epoch] = append(res[epoch], a)
	}

	return res
}

func bucketActionsByDay(actions []types.StakingAction) map[string][]types.StakingAction {
	res := make(map[string][]types.StakingAction)

	for _, a := range actions {
		day := time.Unix(a.BlockTimestamp, 0).UTC().Format("01/02/2006")
		res[day] = append(res[day], a)
	}

	return res
}
func calcSum(actions []types.StakingAction) types.Aggregate {
	var agg types.Aggregate

	for _, a := range actions {
		if a.ActionType == "DEPOSIT" {
			agg.SumDeposits = agg.SumDeposits.Add(scaleDecimals(a.Amount, a.TokenAddress))
		} else {
			agg.SumWithdrawals = agg.SumWithdrawals.Add(scaleDecimals(a.Amount, a.TokenAddress))
		}
	}

	return agg
}

func scaleDecimals(value int64, token string) decimal.Decimal {
	dec := types.Decimals[strings.ToLower(token)]
	valueDec := decimal.New(value, 64)

	return valueDec.DivRound(decimal.New(1, dec), dec)
}

func filterActionsByPool(actions []types.StakingAction, p types.Pool) []types.StakingAction {
	var filtered []types.StakingAction
	for _, a := range actions {
		if strInArrayLowercase(a.TokenAddress, p.Tokens) {
			filtered = append(filtered, a)
		}
	}

	return filtered
}

func strInArrayLowercase(str string, slice []string) bool {
	for _, v := range slice {
		if strings.ToLower(str) == strings.ToLower(v) {
			return true
		}
	}

	return false
}
