package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/utils"
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
}

type juniorRedeem struct {
	JuniorBondAddress string          `json:"juniorBondAddress"`
	UserAddress       string          `json:"userAddress"`
	JuniorBondID      int64           `json:"juniorBondId"`
	SYAddress         string          `json:"smartYieldAddress"`
	TokensIn          decimal.Decimal `json:"tokensIn"`
	MaturesAt         int64           `json:"maturesAt"`
	UnderlyingOut     decimal.Decimal `json:"underlyingOut"`
	BlockTimestamp    int64           `json:"blockTimestamp"`
}

func (a *API) handleSeniorRedeems(c *gin.Context) {
	userAddress := utils.NormalizeAddress(c.Param("address"))

	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	var seniorBondRedeems []seniorRedeem
	rows, err := a.db.Query(`
		select r.sy_address,
			   r.owner_address,
			   r.senior_bond_address,
			   r.fee,
			   r.block_timestamp,
			   r.senior_bond_id,
			   b.underlying_in,
			   b.gain,
			   b.for_days
		from smart_yield_senior_redeem as r
				 inner join smart_yield_senior_buy as b
							on r.senior_bond_address = b.senior_bond_address and r.senior_bond_id = b.senior_bond_id
		where r.owner_address = $1
		offset $2 limit $3
	`, userAddress, offset, limit)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	for rows.Next() {
		var redeem seniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.UserAddress, &redeem.SeniorBondAddress, &redeem.Fee, &redeem.BlockTimestamp, &redeem.SeniorBondID, &redeem.UnderlyingIn, &redeem.Gain, &redeem.ForDays)
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
	err = a.db.QueryRow(`
		select count(*)
		from smart_yield_senior_redeem as r
				 inner join smart_yield_senior_buy as b
							on r.senior_bond_address = b.senior_bond_address and r.senior_bond_id = b.senior_bond_id
		where r.owner_address = $1
	`, userAddress).Scan(&count)

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, seniorBondRedeems, map[string]interface{}{"count": count, "block": block})
}

func (a *API) handleJuniorRedeems(c *gin.Context) {
	userAddress := utils.NormalizeAddress(c.Param("address"))

	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	var juniorBondRedeems []juniorRedeem
	rows, err := a.db.Query(`
		select r.sy_address,
			   r.owner_address,
			   r.junior_bond_address,
			   r.junior_bond_id,
			   r.underlying_out,
			   r.block_timestamp,
			   b.tokens_in,
			   b.matures_at
		from smart_yield_junior_redeem as r
				 inner join smart_yield_junior_buy as b
							on r.junior_bond_address = b.junior_bond_address and r.junior_bond_id = b.junior_bond_id
		where r.owner_address = $1
		offset $2 limit $3
	`, userAddress, offset, limit)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	for rows.Next() {
		var redeem juniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.UserAddress, &redeem.JuniorBondAddress, &redeem.JuniorBondID, &redeem.UnderlyingOut, &redeem.BlockTimestamp, &redeem.TokensIn, &redeem.MaturesAt)
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
		redeem.TokensIn = redeem.TokensIn.DivRound(tenPowDec, int32(p.UnderlyingDecimals))
		redeem.UnderlyingOut = redeem.UnderlyingOut.DivRound(tenPowDec, int32(p.UnderlyingDecimals))

		juniorBondRedeems = append(juniorBondRedeems, redeem)
	}

	var count int64
	err = a.db.QueryRow(`
		select count(*)
		from smart_yield_junior_redeem as r
				 inner join smart_yield_junior_buy as b
							on r.junior_bond_address = b.junior_bond_address and r.junior_bond_id = b.junior_bond_id
		where r.owner_address = $1
	`, userAddress).Scan(&count)

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, juniorBondRedeems, map[string]interface{}{"count": count, "block": block})
}
