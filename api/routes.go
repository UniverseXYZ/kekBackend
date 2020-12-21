package api

func (a *API) setRoutes() {
	explorer := a.engine.Group("/api/explorer")
	explorer.GET("/block/:block", a.BlockHandler)
	explorer.GET("/block-range/:start/:end", a.BlockRangeHandler)
	explorer.GET("/uncle/:hash", a.UncleDetailsHandler)
	explorer.GET("/tx/:txHash", a.TxDetailsHandler)
	explorer.GET("/tx/:txHash/log-entries", a.TxLogEntriesHandler)
	explorer.GET("/search/:query", a.SearchHandler)

	explorer.GET("/account/:address/txs", a.AccountTxsHandler)
	explorer.GET("/account/:address/code", a.AccountCodeHandler)
	explorer.GET("/account/:address/balance", a.AccountBalanceHandler)

	governance := a.engine.Group("/api/governance")
	governance.GET("/proposals/:proposalID", a.ProposalDetailsHandler)
	governance.GET("/proposals/", a.AllProposalHandler)
	governance.GET("/votes/:proposalID", a.VoteDetailsHandler)
}
