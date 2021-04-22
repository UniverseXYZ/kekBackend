package api

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/kekDAO/kekBackend/api/types"
	"github.com/kekDAO/kekBackend/utils"
)

func (a *API) handleStakingActionsList(c *gin.Context) {
	filters := new(Filters)

	userAddress := c.DefaultQuery("userAddress", "all")
	if userAddress != "all" {
		filters.Add("user_address", utils.NormalizeAddress(userAddress))
	}

	actionType := strings.ToUpper(c.DefaultQuery("actionType", "all"))

	if actionType != "ALL" {
		if !checkTxType(actionType) {
			Error(c, errors.New("action type does not exist"))
			return
		}

		filters.Add("action_type", actionType)
	}

	tokenAddress := strings.ToLower(c.DefaultQuery("tokenAddress", "all"))

	if tokenAddress != "all" {
		filters.Add("token_address", tokenAddress)
	}

	limit, err := getQueryLimit(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	page, err := getQueryPage(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	offset := (page - 1) * limit

	query, params := buildQueryWithFilter(`select 
				tx_hash,
				user_address,
				token_address,
				amount ,
				action_type,
				block_timestamp
	from yield_farming_actions
	%s
	order by included_in_block desc, tx_index desc, log_index desc
	%s %s`,
		filters,
		&limit,
		&offset)

	stakingActions, err := a.execQuery(query, params)
	if err != nil {
		Error(c, err)
		return
	}

	query, params = buildQueryWithFilter(`select count(*) from  yield_farming_actions %s %s %s`, filters, nil, nil)

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

	OK(c, stakingActions, map[string]interface{}{
		"block": block,
		"count": count,
	})

}

func checkTxType(action string) bool {
	txType := [2]string{"DEPOSIT", "WITHDRAW"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}

func (a *API) handleStakinsActionsChart(c *gin.Context) {
	tokensAddress := strings.ToLower(c.DefaultQuery("tokenAddress", ""))

	if tokensAddress == "" {
		BadRequest(c, errors.New("tokenAddress required"))
		return
	}

	tokens := strings.Split(tokensAddress, ",")

	for i, token := range tokens {
		t, err := utils.ValidateAccount(token)
		if err != nil {
			BadRequest(c, err)
			return
		}
		tokens[i] = t
	}

	startTsString := strings.ToLower(c.DefaultQuery("start", "-1"))
	startTs, err := validateTs(startTsString)
	if err != nil {
		BadRequest(c, err)
		return
	}

	endTsString := strings.ToLower(c.DefaultQuery("end", "-1"))
	endTs, err := validateTs(endTsString)
	if err != nil {
		BadRequest(c, err)
		return
	}

	scale := strings.ToLower(c.DefaultQuery("scale", "week"))
	if scale != "week" && scale != "day" {
		BadRequest(c, errors.New("Wrong scale"))
		return
	}

	charts := make(map[string]types.Chart)

	for _, token := range tokens {
		rows, err := a.db.Query(`select * from yf_stats_by_token($1,$2,$3,$4) order by point`, token, startTs.Unix(), endTs.Unix(), scale)
		if err != nil {
			Error(c, err)
			return
		}

		chart, err := getChart(rows)
		if err != nil {
			Error(c, err)
			return
		}
		charts[token] = *chart

	}

	OK(c, charts)
	return
}

func getChart(rows *sql.Rows) (*types.Chart, error) {
	x := make(types.Chart)

	for rows.Next() {
		var t time.Time
		var a types.Aggregate
		err := rows.Scan(&t, &a.SumDeposits, &a.SumWithdrawals)
		if err != nil {
			return nil, err
		}

		x[t] = a
	}
	return &x, nil
}

func validateTs(ts string) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(ts, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "invalid timestamp")
	}

	if timestamp == -1 {
		return nil, errors.New("timestamp required")
	}

	tm := time.Unix(timestamp, 0)

	return &tm, nil
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
			user            string
			token           string
			amount          decimal.Decimal
			transactionHash string
			actionType      string
			blockTimestamp  int64
		)
		err := rows.Scan(&transactionHash, &user, &token, &amount, &actionType, &blockTimestamp)
		if err != nil {
			return nil, err
		}

		stakingAction := types.StakingAction{
			UserAddress:     user,
			TokenAddress:    token,
			Amount:          amount,
			TransactionHash: transactionHash,
			ActionType:      actionType,
			BlockTimestamp:  blockTimestamp,
		}
		stakingActions = append(stakingActions, stakingAction)

	}
	return stakingActions, nil
}
