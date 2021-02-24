package api

import (
	"database/sql"

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
}

type juniorRedeem struct {
	JuniorBondAddress string `json:"juniorBondAddress"`
	JuniorBondID      int64  `json:"juniorBondId"`
	SYAddress         string `json:"smartYieldAddress"`
	TokensIn          int64  `json:"tokensIn"`
	MaturesAt         int64  `json:"maturesAt"`
	UnderlyingOut     int64  `json:"underlyingOut"`
}

func (a *API) handleSeniorRedeems(c *gin.Context) {
	userAddress := c.Param("user")
	var seniorBondRedeems []seniorRedeem
	rows, err := a.db.Query(`select smart_yield_senior_buy.sy_address , smart_yield_senior_buy.buyer_address,smart_yield_senior_buy.senior_bond_id,smart_yield_senior_buy.underlying_in,
									smart_yield_senior_buy.gain ,smart_yield_senior_buy.for_days,
									smart_yield_senior_redeem.fee 
								from smart_yield_senior_buy
								inner join smart_yield_senior_redeem
								    on smart_yield_senior_redeem.owner_address = smart_yield_senior_buy.buyer_address
								           and smart_yield_senior_redeem.senior_bond_id = smart_yield_senior_buy.senior_bond_id
								where smart_yield_senior_redeem.owner_address = $1`, userAddress)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	for rows.Next() {
		var redeem seniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.SeniorBondAddress, &redeem.SeniorBondID, &redeem.UnderlyingIn, &redeem.Gain, &redeem.ForDays, &redeem.Fee)
		if err != nil {
			Error(c, err)
			return
		}

		seniorBondRedeems = append(seniorBondRedeems, redeem)
	}

	OK(c, seniorBondRedeems)
}

func (a *API) handleJuniorRedeems(c *gin.Context) {
	userAddress := c.Param("user")
	var juniorBondRedeems []juniorRedeem
	rows, err := a.db.Query(`select smart_yield_junior_buy.sy_address,smart_yield_junior_buy.buyer_address,smart_yield_junior_buy.junior_bond_id,
									smart_yield_junior_buy.tokens_in,smart_yield_junior_buy.matures_at,
									smart_yield_junior_redeem.underlying_out
									from smart_yield_junior_buy
									inner join smart_yield_junior_redeem 
									    on smart_yield_junior_redeem.owner_address = smart_yield_junior_buy.buyer_address
								           and smart_yield_junior_redeem.junior_bond_id = smart_yield_junior_buy.junior_bond_id
								where smart_yield_junior_redeem.owner_address = $1`, userAddress)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	for rows.Next() {
		var redeem juniorRedeem
		err := rows.Scan(&redeem.SYAddress, &redeem.JuniorBondAddress, &redeem.JuniorBondID, &redeem.TokensIn, &redeem.MaturesAt, &redeem.UnderlyingOut)
		if err != nil {
			Error(c, err)
			return
		}

		juniorBondRedeems = append(juniorBondRedeems, redeem)
	}

	OK(c, juniorBondRedeems)
}
