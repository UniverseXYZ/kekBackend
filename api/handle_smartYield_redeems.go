package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
)

type seniorRedeem struct {
	SeniorBondAddress string          `json:"seniorBondAddress"`
	UserAddress       string          `json:"userAddress"`
	SeniorBondID      int64           `json:"seniorBondId"`
	SYAddress         string          `json:"smartYieldAddress"`
	Fee               decimal.Decimal `json:"fee"`
	UnderlyingIn      decimal.Decimal `json:"underlyingIn"`
	Gain              decimal.Decimal `json:"gain"`
	ForDays           int64           `json:"forDays"`
	BlockTimestamp    int64           `json:"blockTimestamp"`
	TxHash            string          `json:"transactionHash"`
}

func (a *API) handleSeniorRedeems(c *gin.Context) {
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
	filters.Add("owner_address", user)

	originator := strings.ToLower(c.DefaultQuery("originator", "all"))
	if originator != "all" {
		if !IsSupportedOriginator(originator) {
			BadRequest(c, errors.New("invalid originator parameter"))
			return
		}

		filters.Add("( select p.protocol_id from smart_yield_pools as p where p.sy_address = r.sy_address )", originator)
	}

	token := strings.ToLower(c.DefaultQuery("token", "all"))
	if token != "all" {
		if state.PoolByUnderlyingAddress(token) == nil {
			BadRequest(c, errors.New("invalid token parameter"))
			return
		}

		filters.Add("( select p.underlying_address from smart_yield_pools as p where p.sy_address = r.sy_address )", token)
	}

	query, params := buildQueryWithFilter(`
		select r.sy_address,
			   r.owner_address,
			   r.senior_bond_address,
			   r.fee,
			   r.block_timestamp,
			   r.senior_bond_id,
			   r.tx_hash,
			   b.underlying_in,
			   b.gain,
			   b.for_days
		from smart_yield_senior_redeem as r
				 inner join smart_yield_senior_buy as b
							on r.senior_bond_address = b.senior_bond_address and r.senior_bond_id = b.senior_bond_id
		where %s
		order by r.included_in_block desc, r.tx_index desc, r.log_index desc
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

	var seniorBondRedeems []seniorRedeem
	for rows.Next() {
		var redeem seniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.UserAddress, &redeem.SeniorBondAddress, &redeem.Fee, &redeem.BlockTimestamp, &redeem.SeniorBondID, &redeem.TxHash, &redeem.UnderlyingIn, &redeem.Gain, &redeem.ForDays)
		if err != nil {
			Error(c, err)
			return
		}

		p := state.PoolBySmartYieldAddress(redeem.SYAddress)
		if p == nil {
			Error(c, errors.New("could not find pool in state"))
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

		redeem.UnderlyingIn = redeem.UnderlyingIn.DivRound(tenPowDec, int32(p.UnderlyingDecimals))
		redeem.Fee = redeem.Fee.DivRound(tenPowDec, int32(p.UnderlyingDecimals))
		redeem.Gain = redeem.Gain.DivRound(tenPowDec, int32(p.UnderlyingDecimals))

		seniorBondRedeems = append(seniorBondRedeems, redeem)
	}

	var count int64
	query, params = buildQueryWithFilter(`
		select count(*)
		from smart_yield_senior_redeem as r
				 inner join smart_yield_senior_buy as b
							on r.senior_bond_address = b.senior_bond_address and r.senior_bond_id = b.senior_bond_id
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

	OK(c, seniorBondRedeems, map[string]interface{}{"count": count, "block": block})
}
