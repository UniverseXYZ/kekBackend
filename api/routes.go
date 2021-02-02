package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) setRoutes() {
	a.engine.GET("/health", func(context *gin.Context) {
		err := a.db.Ping()
		if err != nil {
			context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": http.StatusInternalServerError,
				"data":   "NOT OK",
			})

			return
		}

		OK(context, "OK")
	})

	governance := a.engine.Group("/api/governance")
	governance.GET("/proposals", a.AllProposalHandler)
	governance.GET("/proposals/:proposalID", a.ProposalDetailsHandler)
	governance.GET("/proposals/:proposalID/votes", a.VotesHandler)
	governance.GET("/proposals/:proposalID/events", a.handleProposalEvents)
	governance.GET("/overview", a.BondOverview)
	governance.GET("/voters", a.handleVoters)
	governance.GET("/abrogation-proposals", a.AllAbrogationProposals)
	governance.GET("/abrogation-proposals/:proposalID", a.AbrogationProposalDetailsHandler)
	governance.GET("/abrogation-proposals/:proposalID/votes", a.AbrogationVotesHandler)

	yieldfarming := a.engine.Group("/api/yieldfarming")
	yieldfarming.GET("/staking-actions/list", a.handleStakingActions)
}
