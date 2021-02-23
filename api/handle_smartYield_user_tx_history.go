package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) handleSYUserTransactionHistory(c *gin.Context) {
	user := c.Param("address")

	userAddress, err := utils.ValidateAccount(user)
	if err != nil {
		BadRequest(c, errors.New("invalid user address"))
		return
	}

	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	rows, err := a.db.Query(` select protocol_id, sy_address, underlying_token_address, amount, tranche, transaction_type, tx_hash, block_timestamp, included_in_block from smart_yield_transaction_history where user_address = $1 order by included_in_block desc, tx_index desc, log_index desc offset $2 limit $3;`, userAddress, offset, limit)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var history []types.SYUserHistory
	for rows.Next() {
		var h types.SYUserHistory

		err := rows.Scan(&h.ProtocolId, &h.Pool, &h.UnderlyingTokenAddress, &h.Amount, &h.Tranche, &h.TransactionType, &h.TransactionHash, &h.BlockTimestamp, &h.BlockNumber)
		if err != nil {
			Error(c, err)
			return
		}

		p := state.PoolBySmartYieldAddress(h.Pool)
		h.Amount = h.Amount.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals)), int32(p.UnderlyingDecimals))

		history = append(history, h)
	}

	var count int64
	err = a.db.QueryRow(`select count(*) from smart_yield_transaction_history where user_address = $1`, userAddress).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, history, map[string]interface{}{"count": count, "block": block})
}
