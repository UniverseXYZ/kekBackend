package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
)

type JuniorPastPosition struct {
	ProtocolId             string `json:"protocolId"`
	SmartYieldAddress      string `json:"smartYieldAddress"`
	UnderlyingTokenAddress string `json:"underlyingTokenAddress"`
	UnderlyingTokenSymbol  string `json:"underlyingTokenSymbol"`

	TokensIn        decimal.Decimal `json:"tokensIn"`
	UnderlyingOut   decimal.Decimal `json:"underlyingOut"`
	Forfeits        decimal.Decimal `json:"forfeits"`
	TransactionType string          `json:"transactionType"`

	BlockTimestamp int64  `json:"blockTimestamp"`
	TxHash         string `json:"transactionHash"`
}

func (a *API) handleJuniorPastPositions(c *gin.Context) {
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

	filters := new(Filters)
	filters.Add("user_address", user)

	originator := strings.ToLower(c.DefaultQuery("originator", "all"))
	if originator != "all" {
		if !IsSupportedOriginator(originator) {
			BadRequest(c, errors.New("invalid originator parameter"))
			return
		}

		filters.Add("protocol_id", originator)
	}

	token := strings.ToLower(c.DefaultQuery("token", "all"))
	if token != "all" {
		if state.PoolByUnderlyingAddress(token) == nil {
			BadRequest(c, errors.New("invalid token parameter"))
			return
		}

		filters.Add("underlying_token_address", token)
	}

	query, params := buildQueryWithFilter(`
		select h.protocol_id as protocol_id,
			   h.sy_address,
			   underlying_token_address,
			   ( select underlying_symbol from smart_yield_pools p where h.sy_address = p.sy_address ) as underlying_token_symbol,
			   coalesce(
				   (select tokens_in from smart_yield_token_sell s where s.tx_hash = h.tx_hash and s.log_index = h.log_index ),
				   (select tokens_in from smart_yield_junior_redeem r inner join smart_yield_junior_buy syjb on r.junior_bond_address = syjb.junior_bond_address and r.junior_bond_id = syjb.junior_bond_id where r.tx_hash = h.tx_hash and r.log_index = h.log_index )
			   ) as tokens_in,
			   coalesce(
				   ( select underlying_out from smart_yield_token_sell s where s.tx_hash = h.tx_hash and s.log_index = h.log_index ),
				   ( select underlying_out from smart_yield_junior_redeem r where r.tx_hash = h.tx_hash and r.log_index = h.log_index )
			   ) as underlying_out,
			   coalesce(
				   ( select forfeits from smart_yield_token_sell s where s.tx_hash = h.tx_hash and s.log_index = h.log_index),
				   0
			   ) as forfeits,
			   transaction_type,
			   tx_hash,
			   block_timestamp
		from smart_yield_transaction_history h
		where %s and transaction_type = ANY(ARRAY['JUNIOR_INSTANT_WITHDRAW'::sy_tx_history_tx_type, 'JUNIOR_REDEEM'::sy_tx_history_tx_type])
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

	var positions []JuniorPastPosition
	for rows.Next() {
		var p JuniorPastPosition

		err := rows.Scan(&p.ProtocolId, &p.SmartYieldAddress, &p.UnderlyingTokenAddress, &p.UnderlyingTokenSymbol, &p.TokensIn, &p.UnderlyingOut, &p.Forfeits, &p.TransactionType, &p.TxHash, &p.BlockTimestamp)
		if err != nil {
			Error(c, err)
			return
		}

		pool := state.PoolBySmartYieldAddress(p.SmartYieldAddress)
		if pool == nil {
			Error(c, errors.New("could not find pool in state"))
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(pool.UnderlyingDecimals))

		p.TokensIn = p.TokensIn.DivRound(tenPowDec, int32(pool.UnderlyingDecimals))
		p.UnderlyingOut = p.UnderlyingOut.DivRound(tenPowDec, int32(pool.UnderlyingDecimals))
		p.Forfeits = p.Forfeits.DivRound(tenPowDec, int32(pool.UnderlyingDecimals))

		positions = append(positions, p)
	}

	var count int64
	query, params = buildQueryWithFilter(`
		select count(*)
		from smart_yield_transaction_history as h 
		where %s and transaction_type = ANY(ARRAY['JUNIOR_INSTANT_WITHDRAW'::sy_tx_history_tx_type, 'JUNIOR_REDEEM'::sy_tx_history_tx_type])
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

	OK(c, positions, map[string]interface{}{"count": count, "block": block})
}
