package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/utils"
)

type seniorRedeem struct {
	SeniorBondAddress string `json:"seniorBondAddress"`
	UserAddress       string `json:"userAddress"`
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
	UserAddress       string `json:"userAddress"`
	JuniorBondID      int64  `json:"juniorBondId"`
	SYAddress         string `json:"smartYieldAddress"`
	TokensIn          int64  `json:"tokensIn"`
	MaturesAt         int64  `json:"maturesAt"`
	UnderlyingOut     int64  `json:"underlyingOut"`
	BlockTimestamp    int64  `json:"blockTimestamp"`
}

func (a *API) handleSeniorRedeems(c *gin.Context) {
	userAddress := utils.NormalizeAddress(c.Param("address"))
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
									on r.senior_bond_address = b.senior_bond_address
										and r.senior_bond_id = b.senior_bond_id
				where r.owner_address = $1 
				`, userAddress)

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

		seniorBondRedeems = append(seniorBondRedeems, redeem)
	}

	OK(c, seniorBondRedeems)
}

func (a *API) handleJuniorRedeems(c *gin.Context) {
	userAddress := utils.NormalizeAddress(c.Param("address"))
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
									on r.junior_bond_address = b.junior_bond_address
										and r.junior_bond_id = b.junior_bond_id
				where r.owner_address = $1
				`, userAddress)

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

		juniorBondRedeems = append(juniorBondRedeems, redeem)
	}

	OK(c, juniorBondRedeems)
}
