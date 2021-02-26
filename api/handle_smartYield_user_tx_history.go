package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYield"
	"github.com/barnbridge/barnbridge-backend/state"
)

func (a *API) handleSYUserTransactionHistory(c *gin.Context) {
	user, err := getQueryAddress(c, "address")
	if err != nil {
		BadRequest(c, err)
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

	filters := make(map[string]interface{})
	filters["user_address"] = user

	originator := strings.ToLower(c.DefaultQuery("originator", "all"))
	if originator != "all" {
		if !IsSupportedOriginator(originator) {
			BadRequest(c, errors.New("invalid originator parameter"))
			return
		}

		filters["protocol_id"] = originator
	}

	token := strings.ToLower(c.DefaultQuery("token", "all"))
	if token != "all" {
		if state.PoolByUnderlyingAddress(token) == nil {
			BadRequest(c, errors.New("invalid token parameter"))
			return
		}

		filters["underlying_token_address"] = token
	}

	txType := strings.ToUpper(c.DefaultQuery("transactionType", "all"))
	if txType != "ALL" {
		if !IsSupportedTxType(txType) {
			BadRequest(c, errors.New("invalid transactionType parameter"))
			return
		}

		filters["transaction_type"] = txType
	}

	query, params := buildQueryWithFilter(`
		select h.protocol_id,
			   h.sy_address,
			   underlying_token_address,
               (select underlying_symbol from smart_yield_pools p where h.sy_address = p.sy_address) as underlying_token_symbol, 
			   amount,
			   tranche,
			   transaction_type,
			   tx_hash,
			   block_timestamp,
			   included_in_block
		from smart_yield_transaction_history h
		where %s
		order by included_in_block desc, tx_index desc, log_index desc
		%s %s;
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

	var history []types.SYUserHistory
	for rows.Next() {
		var h types.SYUserHistory

		err := rows.Scan(&h.ProtocolId, &h.Pool, &h.UnderlyingTokenAddress, &h.UnderlyingTokenSymbol, &h.Amount, &h.Tranche, &h.TransactionType, &h.TransactionHash, &h.BlockTimestamp, &h.BlockNumber)
		if err != nil {
			Error(c, err)
			return
		}

		p := state.PoolBySmartYieldAddress(h.Pool)
		h.Amount = h.Amount.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals)), int32(p.UnderlyingDecimals))

		history = append(history, h)
	}

	query, params = buildQueryWithFilter(`
		select count(*)
		from smart_yield_transaction_history
		where %s
		%s %s;
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

	OK(c, history, map[string]interface{}{"count": count, "block": block})
}

func IsSupportedOriginator(originator string) bool {
	switch strings.ToLower(originator) {
	case "compound/v2":
		return true
	}

	return false
}

func IsSupportedTxType(t string) bool {
	switch smartYield.TxType(strings.ToUpper(t)) {
	case smartYield.JuniorDeposit, smartYield.JuniorInstantWithdraw, smartYield.JuniorRegularWithdraw, smartYield.JuniorRedeem, smartYield.SeniorDeposit, smartYield.SeniorRedeem, smartYield.JtokenSend, smartYield.JtokenReceive, smartYield.JbondSend, smartYield.JbondReceive, smartYield.SbondSend, smartYield.SbondReceive:
		return true
	}

	return false
}
