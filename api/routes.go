package api

func (a *API) setRoutes() {
	governance := a.engine.Group("/api/governance")
	governance.GET("/proposals/:proposalID", a.ProposalDetailsHandler)
	governance.GET("/proposals/", a.AllProposalHandler)
	governance.GET("/votes/:proposalID", a.VotesHandler)
}
