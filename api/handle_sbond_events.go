package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleSBondEvents(c *gin.Context) {
	userAddress := strings.ToLower(c.DefaultQuery("user", ""))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	trades, err := a.getAllJBondEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]interface{}{"Buys": trades.Buys, "Redeems": trades.Buys, "Transfers": trades.Transfers})
}

func (a *API) getAllSBondEvents(userAddress string, offset string, limit string) (*types.SeniorBondTrades, error) {
	var trades types.SeniorBondTrades
	var err error

	trades.Buys, err = a.getAllSBuyEvents(userAddress, offset, limit)
	if err != nil {
		return nil, err
	}

	trades.Redeems, err = a.getAllSRedeemEvent(userAddress, offset, limit)
	if err != nil {
		return nil, err
	}

	trades.Transfers, err = a.getERC721Transfer("senior", userAddress, offset, limit)

	return &trades, nil
}

func (a *API) getAllSBuyEvents(userAddress string, offset string, limit string) ([]types.SeniorBondBuyTrade, error) {
	query := `
			select  sy_address,
					buyer_address,
					senior_bond_id,
					underlying_in,
					gain,
			        for_days,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from smart_yield_senior_buy
		where 1=1 %s order by block_timestamp desc
		offset $1 
		limit $2 `

	var parameters = []interface{}{offset, limit}

	if userAddress != "" {
		parameters = append(parameters, userAddress)
		userFilter := fmt.Sprintf("and buyer_address = $%d", len(parameters))
		query = fmt.Sprintf(query, userFilter)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := a.db.Query(query, parameters...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	var sBondBuys []types.SeniorBondBuyTrade
	for rows.Next() {
		var t types.SeniorBondBuyTrade
		err := rows.Scan(&t.SYAddress, &t.BuyerAddress, &t.SeniorBondID, &t.UnderlyingIn, &t.Gain, &t.ForDays, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		sBondBuys = append(sBondBuys, t)
	}

	return sBondBuys, nil
}

func (a *API) getAllSRedeemEvent(userAddress string, offset string, limit string) ([]types.SeniorBondRedeemTrade, error) {
	query := `
			select  sy_address,
					owner_address,
					senior_bond_id,
					fee,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from smart_yield_senior_redeem
		where 1=1 %s order by block_timestamp desc
		offset $1 
		limit $2 `

	var parameters = []interface{}{offset, limit}

	if userAddress != "" {
		parameters = append(parameters, userAddress)
		userFilter := fmt.Sprintf("and owner_address = $%d", len(parameters))
		query = fmt.Sprintf(query, userFilter)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := a.db.Query(query, parameters...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	var SBondRedeems []types.SeniorBondRedeemTrade

	for rows.Next() {
		var t types.SeniorBondRedeemTrade
		err := rows.Scan(&t.SYAddress, &t.OwnerAddress, &t.SeniorBondID, &t.Fee, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		SBondRedeems = append(SBondRedeems, t)
	}

	return SBondRedeems, nil
}
