package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleSYTokenEvents(c *gin.Context) {
	userAddress := strings.ToLower(c.DefaultQuery("user", ""))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	trades, err := a.getAllSYTokenEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]interface{}{"Buys": trades.Buys, "Redeems": trades.Sells})
}

func (a *API) getAllSYTokenEvents(userAddress string, offset string, limit string) (*types.SYTokenTrades, error) {
	var trades types.SYTokenTrades
	var err error

	trades.Buys, err = a.getSYTokenBuyEvents(userAddress, offset, limit)
	if err != nil {
		return nil, err
	}

	trades.Sells, err = a.getAllSYTokenSellEvents(userAddress, offset, limit)
	if err != nil {
		return nil, err
	}

	return &trades, nil
}

func (a *API) getSYTokenBuyEvents(userAddress string, offset string, limit string) ([]types.TokenBuyTrade, error) {
	query := `
			select  sy_address,
					buyer_address,
					underlying_in,
					tokens_out,
			        fee,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from smart_yield_token_buy
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

	var buys []types.TokenBuyTrade
	for rows.Next() {
		var t types.TokenBuyTrade
		err := rows.Scan(&t.SYAddress, &t.BuyerAddress, &t.UnderlyingIn, &t.TokensOut, &t.Fee, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		buys = append(buys, t)
	}

	return buys, nil
}

func (a *API) getAllSYTokenSellEvents(userAddress string, offset string, limit string) ([]types.TokenSellTrade, error) {
	query := `
			select  sy_address,
					seller_address,
					tokens_in,
					underlying_out,
			       	forfeits,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from smart_yield_token_sell
		where 1=1 %s order by block_timestamp desc
		offset $1 
		limit $2 `

	var parameters = []interface{}{offset, limit}

	if userAddress != "" {
		parameters = append(parameters, userAddress)
		userFilter := fmt.Sprintf("and seller_address = $%d", len(parameters))
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

	var tokenSells []types.TokenSellTrade

	for rows.Next() {
		var t types.TokenSellTrade
		err := rows.Scan(&t.SYAddress, &t.SellerAddress, &t.TokensIn, &t.UnderlyingOut, &t.Forfeits, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		tokenSells = append(tokenSells, t)
	}

	return tokenSells, nil
}
