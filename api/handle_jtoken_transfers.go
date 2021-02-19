package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handleJTokenTransfer(c *gin.Context) {
	userAddress := strings.ToLower(c.DefaultQuery("user", ""))
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	transfers, err := a.getAllJTokenTransfer(userAddress, offset, limit)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]interface{}{"Transfers": transfers})
}

func (a *API) getAllJTokenTransfer(userAddress string, offset string, limit string) ([]types.Transfer, error) {
	query := `
			select  sy_address,
					sender,
					receiver,
					value,
			       
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from jtoken_transfers
		where 1=1 %s order by block_timestamp desc
		offset $1 
		limit $2 `

	var parameters = []interface{}{offset, limit}

	if userAddress != "" {
		parameters = append(parameters, userAddress)
		userFilter := fmt.Sprintf("and sender = $%d or receiver = $%d", len(parameters), len(parameters))
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

	var transfers []types.Transfer
	for rows.Next() {
		var t types.Transfer
		err := rows.Scan(&t.SYAddress, &t.From, &t.To, &t.Value, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}
