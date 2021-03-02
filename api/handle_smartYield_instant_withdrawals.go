package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
)

type JuniorInstantWithdrawal struct {
	SmartYieldAddress      string          `json:"smartYieldAddress"`
	TokensIn               decimal.Decimal `json:"tokensIn"`
	UnderlyingOut          decimal.Decimal `json:"underlyingOut"`
	Forfeits               decimal.Decimal `json:"forfeits"`
	BlockTimestamp         int64           `json:"blockTimestamp"`
	TxHash                 string          `json:"transactionHash"`
	UnderlyingTokenAddress string          `json:"underlyingTokenAddress"`
	ProtocolId             string          `json:"protocolId"`
}

func (a *API) handleJuniorInstantWithdrawals(c *gin.Context) {
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
	filters["seller_address"] = user

	originator := strings.ToLower(c.DefaultQuery("originator", "all"))
	if originator != "all" {
		if !IsSupportedOriginator(originator) {
			BadRequest(c, errors.New("invalid originator parameter"))
			return
		}

		filters["( select p.protocol_id from smart_yield_pools as p where p.sy_address = r.sy_address )"] = originator
	}

	token := strings.ToLower(c.DefaultQuery("token", "all"))
	if token != "all" {
		if state.PoolByUnderlyingAddress(token) == nil {
			BadRequest(c, errors.New("invalid token parameter"))
			return
		}

		filters["( select p.underlying_address from smart_yield_pools as p where p.sy_address = r.sy_address )"] = token
	}

	query, params := buildQueryWithFilter(`
		select r.sy_address,
			   r.tokens_in,
			   r.underlying_out,
			   r.forfeits,
			   r.block_timestamp,
			   r.tx_hash,
			   ( select p.underlying_address from smart_yield_pools as p where p.sy_address = r.sy_address ) as underlying_token_address,
			   ( select p.protocol_id from smart_yield_pools as p where p.sy_address = r.sy_address ) as protocol_id
		from smart_yield_token_sell as r
		where %s
		order by included_in_block desc, tx_index desc, log_index desc
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

	var instantWithdrawals []JuniorInstantWithdrawal
	for rows.Next() {
		var w JuniorInstantWithdrawal
		err := rows.Scan(&w.SmartYieldAddress, &w.TokensIn, &w.UnderlyingOut, &w.Forfeits, &w.BlockTimestamp, &w.TxHash, &w.UnderlyingTokenAddress, &w.ProtocolId)
		if err != nil {
			Error(c, err)
			return
		}

		p := state.PoolBySmartYieldAddress(w.SmartYieldAddress)
		if p == nil {
			Error(c, errors.New("could not find pool in state"))
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

		w.TokensIn = w.TokensIn.DivRound(tenPowDec, int32(p.UnderlyingDecimals))
		w.UnderlyingOut = w.UnderlyingOut.DivRound(tenPowDec, int32(p.UnderlyingDecimals))
		w.Forfeits = w.Forfeits.DivRound(tenPowDec, int32(p.UnderlyingDecimals))

		instantWithdrawals = append(instantWithdrawals, w)
	}

	var count int64
	query, params = buildQueryWithFilter(`
		select count(*)
		from smart_yield_token_sell as r 
		where %s
		%s %s;
	`,
		filters,
		nil,
		nil,
	)

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

	OK(c, instantWithdrawals, map[string]interface{}{"count": count, "block": block})
}
