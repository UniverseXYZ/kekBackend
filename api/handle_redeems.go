package api

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
)

type seniorRedeem struct {
	SeniorBondAddress string `json:"seniorBondAddress"`
	SeniorBondID      int64  `json:"seniorBondId"`
	SYAddress         string `json:"smartYieldAddress"`
	Fee               int64  `json:"fee"`
	UnderlyingIn      int64  `json:"underlyingIn"`
	Gain              int64  `json:"gain"`
	ForDays           int64  `json:"forDays"`
	BlockTimestamp    int64  `json:"blockTimestamp"`
}

type juniorRedeem struct {
	JuniorBondAddress string `json:"juniorBondAddress"`
	JuniorBondID      int64  `json:"juniorBondId"`
	SYAddress         string `json:"smartYieldAddress"`
	TokensIn          int64  `json:"tokensIn"`
	MaturesAt         int64  `json:"maturesAt"`
	UnderlyingOut     int64  `json:"underlyingOut"`
	BlockTimestamp    int64  `json:"blockTimestamp"`
}

func (a *API) handleSeniorRedeems(c *gin.Context) {
	userAddress := c.Param("user")
	var seniorBondRedeems []seniorRedeem
	rows, err := a.db.Query(`
				select b.sy_address,
					   b.buyer_address,
					   b.senior_bond_id,
					   b.underlying_in,
					   b.gain,
					   b.for_days,
					   r.fee,
					   r.block_timestamp
				from smart_yield_senior_buy as b
						 inner join smart_yield_senior_redeem as r
									on r.senior_bond_address = b.senior_bond_address
										and r.senior_bond_id = b.senior_bond_id
				where (r.owner_address = $1 or b.buyer_address = $1
				)`, userAddress)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	for rows.Next() {
		var redeem seniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.SeniorBondAddress, &redeem.SeniorBondID, &redeem.UnderlyingIn, &redeem.Gain, &redeem.ForDays, &redeem.Fee, &redeem.BlockTimestamp)
		if err != nil {
			Error(c, err)
			return
		}

		seniorBondRedeems = append(seniorBondRedeems, redeem)
	}

	OK(c, seniorBondRedeems)
}

func (a *API) handleJuniorRedeems(c *gin.Context) {
	userAddress := strings.ToLower(c.Param("user"))
	var juniorBondRedeems []juniorRedeem
	rows, err := a.db.Query(`
				select b.sy_address,
					   b.buyer_address,
					   b.junior_bond_id,
					   b.tokens_in,
					   b.matures_at,
					   r.underlying_out,
					   r.block_timestamp
				from smart_yield_junior_buy as b
						 inner join smart_yield_junior_redeem as r
									on r.junior_bond_address = b.junior_bond_address
										and r.junior_bond_id = b.junior_bond_id
				where (r.owner_address = $1 or b.buyer_address = $1
				)`, userAddress)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	for rows.Next() {
		var redeem juniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.JuniorBondAddress, &redeem.JuniorBondID, &redeem.TokensIn, &redeem.MaturesAt, &redeem.UnderlyingOut, &redeem.BlockTimestamp)
		if err != nil {
			Error(c, err)
			return
		}

		juniorBondRedeems = append(juniorBondRedeems, redeem)
	}

	OK(c, juniorBondRedeems)
}
