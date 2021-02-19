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

	jTrades, err := a.getAllJBondEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	sTrades, err := a.getAllSBondEvents(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]map[string]interface{}{"JuniorBond": {"Buys": jTrades.Buys, "Redeems": jTrades.Buys, "Transfers": jTrades.Transfers}, "SeniorBond": {"Buys": sTrades.Buys, "Redeems": sTrades.Buys, "Transfers": sTrades.Transfers}})
}
