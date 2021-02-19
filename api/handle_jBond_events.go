package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleJBondEvents(c *gin.Context) {
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

func (a *API) getAllJBondEvents(userAddress string, offset string, limit string) (*types.JuniorBondTrades, error) {
	var trades types.JuniorBondTrades
	var err error

	trades.Buys, err = a.getAllJBuyEvents(userAddress, offset, limit)
	if err != nil {
		return nil, err
	}

	trades.Redeems, err = a.getAllJRedeemEvent(userAddress, offset, limit)
	if err != nil {
		return nil, err
	}

	trades.Transfers, err = a.getERC721Transfer("junior", userAddress, offset, limit)

	return &trades, nil

}

func (a *API) getAllJBuyEvents(userAddress string, offset string, limit string) ([]types.JuniorBondBuyTrade, error) {
	query := `
			select  sy_address,
					buyer_address,
					junior_bond_id,
					tokens_in,
					matures_at,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from smart_yield_junior_buy
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

	var jBondBuys []types.JuniorBondBuyTrade
	for rows.Next() {
		var t types.JuniorBondBuyTrade
		err := rows.Scan(&t.SYAddress, &t.BuyerAddress, &t.JuniorBondID, &t.TokensIn, &t.MaturesAt, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		jBondBuys = append(jBondBuys, t)
	}

	return jBondBuys, nil
}

func (a *API) getAllJRedeemEvent(userAddress string, offset string, limit string) ([]types.JuniorBondRedeemTrade, error) {
	query := `
			select  sy_address,
					owner_address,
					junior_bond_id,
					underlying_out,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from smart_yield_junior_redeem
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

	var jBondRedeems []types.JuniorBondRedeemTrade

	for rows.Next() {
		var t types.JuniorBondRedeemTrade
		err := rows.Scan(&t.SYAddress, &t.OwnerAddress, &t.JuniorBondID, &t.UnderlyingOut, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		jBondRedeems = append(jBondRedeems, t)
	}

	return jBondRedeems, nil
}
