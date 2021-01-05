package api

func (a *API) setRoutes() {
	governance := a.engine.Group("/api/governance")
	governance.GET("/proposals", a.AllProposalHandler)
	governance.GET("/proposals/:proposalID", a.ProposalDetailsHandler)
	governance.GET("/proposals/:proposalID/votes", a.VotesHandler)
	governance.GET("/overview", a.BondOverview)
	governance.GET("/voters", a.handleVoters)
	governance.GET("/cancellation-proposals", a.AllCancellationProposals)
	governance.GET("/cancellation-proposals/:proposalID", a.CancellationProposalDetailsHandler)
	governance.GET("/cancellation-proposals/:proposalID/votes", a.CancellationVotesHandler)
}
