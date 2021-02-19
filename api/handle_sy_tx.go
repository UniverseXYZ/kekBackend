package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *API) handleSYTxs(c *gin.Context) {
	userAddress := strings.ToLower(c.DefaultQuery("user", ""))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	juniorTrades, err := a.getAllJBondEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	seniorTrades, err := a.getAllSBondEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	syTokenTrades, err := a.getAllSYTokenEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	jTokenTransfers, err := a.getAllJTokenTransfer(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]map[string]interface{}{
		"JuniorBond": {"Buys": juniorTrades.Buys, "Redeems": juniorTrades.Redeems, "Transfers": juniorTrades.Transfers},
		"SeniorBond": {"Buys": seniorTrades.Buys, "Redeems": seniorTrades.Redeems, "Transfers": seniorTrades.Transfers},
		"SYTokens":   {"Buys": syTokenTrades.Buys, "Sells": syTokenTrades.Sells},
		"JToken:":    {"Transfers": jTokenTransfers}})
}
