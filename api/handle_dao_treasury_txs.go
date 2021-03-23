package api

import (
	"database/sql"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type TreasuryTx struct {
	AccountAddress string `json:"accountAddress"`
	AccountLabel   string `json:"accountLabel"`

	CounterpartyAddress string `json:"counterpartyAddress"`
	CounterpartyLabel   string `json:"counterpartyLabel"`

	Amount               decimal.Decimal `json:"amount"`
	TransactionDirection string          `json:"transactionDirection"`
	TokenAddress         string          `json:"tokenAddress"`
	TokenSymbol          string          `json:"tokenSymbol"`
	TokenDecimals        int64           `json:"tokenDecimals"`

	TransactionHash string `json:"transactionHash"`
	BlockTimestamp  int64  `json:"blockTimestamp"`
	BlockNumber     int64  `json:"blockNumber"`
}

func (a *API) handleDaoTxs(c *gin.Context) {
	treasuryAddress := strings.ToLower(c.DefaultQuery("address", ""))
	if treasuryAddress == "" {
		BadRequest(c, errors.New("Address could not be null"))
		return
	}

	limit, err := getQueryLimit(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	page, err := getQueryPage(c)
	if err != nil {
		BadRequest(c, err)
		return
	}

	offset := (page - 1) * limit
	filters := new(Filters)
	filters.Add("t.account", treasuryAddress)

	tokenAddress := strings.ToLower(c.DefaultQuery("tokenAddress", "all"))

	if tokenAddress != "all" {
		filters.Add("token_address", tokenAddress)
	}

	txDirection := strings.ToUpper(c.DefaultQuery("transactionDirection", "all"))

	if txDirection != "ALL" {
		filters.Add("tx_direction", txDirection)
	}

	query, params := buildQueryWithFilter(`
			select t.token_address,
				   t.account,
				   t.counterparty,
				   t.amount,
				   t.tx_direction,
				   t.block_timestamp,
				   t.included_in_block,
				   e20t.symbol,
				   e20t.decimals,
			from account_erc20_transfers as t
					 inner join erc20_tokens e20t
								on t.token_address = e20t.token_address;
			where %s
			%s %s;`, filters, &limit, &offset)

	spew.Dump(query)
	rows, err := a.db.Query(query, params...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var transactions []TreasuryTx

	for rows.Next() {
		var t TreasuryTx

		err := rows.Scan(&t.TokenAddress, &t.AccountAddress, &t.CounterpartyAddress, &t.Amount, &t.TransactionDirection, &t.BlockTimestamp, &t.BlockNumber, &t.TokenSymbol, &t.TokenDecimals, &t.AccountLabel, &t.CounterpartyLabel)
		if err != nil {
			Error(c, err)
			return
		}

		transactions = append(transactions, t)
	}

	OK(c, transactions)
}
