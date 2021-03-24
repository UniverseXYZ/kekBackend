package api

import (
	"database/sql"
	"strings"

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

type Token struct {
	TokenAddress  string `json:"tokenAddress"`
	TokenSymbol   string `json:"tokenSymbol"`
	TokenDecimals int64  `json:"tokenDecimals"`
}

func (a *API) handleTreasuryTxs(c *gin.Context) {
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
		filters.Add("t.token_address", tokenAddress)
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
				   t.tx_hash,
				   e20t.symbol,
				   e20t.decimals,
				   coalesce((select label from labels where address=t.account), '') as accountLabel,
				   coalesce((select label from labels where address=t.counterparty), '') as counterpartyLabel
			from account_erc20_transfers as t
					 inner join erc20_tokens e20t
								on t.token_address = e20t.token_address
			%s order by included_in_block desc,t.tx_index desc,t.log_index desc
			%s %s;`, filters, &limit, &offset)

	rows, err := a.db.Query(query, params...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var transactions []TreasuryTx

	for rows.Next() {
		var t TreasuryTx

		err := rows.Scan(&t.TokenAddress, &t.AccountAddress, &t.CounterpartyAddress, &t.Amount, &t.TransactionDirection, &t.BlockTimestamp, &t.BlockNumber, &t.TransactionHash, &t.TokenSymbol, &t.TokenDecimals, &t.AccountLabel, &t.CounterpartyLabel)
		if err != nil {
			Error(c, err)
			return
		}
		t.Amount = t.Amount.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(t.TokenDecimals)), int32(t.TokenDecimals))
		transactions = append(transactions, t)
	}

	query, params = buildQueryWithFilter(`
			 select count(*) 
			 from account_erc20_transfers t
			 %s 
			 %s %s`,
		filters,
		nil,
		nil)

	var count int64
	err = a.db.QueryRow(query, params...).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, transactions, map[string]interface{}{"count": count, "block": block})
}

func (a *API) handleTreasuryTokens(c *gin.Context) {
	treasuryAddress := strings.ToLower(c.DefaultQuery("address", ""))
	if treasuryAddress == "" {
		BadRequest(c, errors.New("Address could not be null"))
		return
	}

	rows, err := a.db.Query(`
				select distinct 
				                transfers.token_address,
				                tokens.symbol, 
				                tokens.decimals
				from account_erc20_transfers transfers
						 inner join erc20_tokens tokens 
						     on transfers.token_address = tokens.token_address
				where account = $1;`, treasuryAddress)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	var tokens []Token

	for rows.Next() {
		var t Token
		err := rows.Scan(&t.TokenAddress, &t.TokenSymbol, &t.TokenDecimals)
		if err != nil {
			Error(c, err)
			return
		}

		tokens = append(tokens, t)
	}

	OK(c, tokens)
}
