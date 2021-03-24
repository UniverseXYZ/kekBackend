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
	governance.GET("/treasury/transactions", a.handleTreasuryTxs)
	governance.GET("/treasury/tokens", a.handleTreasuryTokens)

	smartYield := a.engine.Group("/api/smartyield")
	smartYield.GET("/pools", a.handlePools)
	smartYield.GET("/pools/:address", a.handlePoolDetails)
	smartYield.GET("/pools/:address/apy", a.handlePoolAPYTrend)
	smartYield.GET("/users/:address/history", a.handleSYUserTransactionHistory)
	smartYield.GET("/users/:address/redeems/senior", a.handleSeniorRedeems)
	smartYield.GET("/users/:address/junior-past-positions", a.handleJuniorPastPositions)
	smartYield.GET("/users/:address/portfolio-value", a.handleSYUserPortfolioValue)
	smartYield.GET("/users/:address/portfolio-value/junior", a.handleSYUserJuniorPortfolioValue)
	smartYield.GET("/users/:address/portfolio-value/senior", a.handleSYUserSeniorPortfolioValue)
	smartYield.GET("/rewards/pools", a.handleRewardPools)
	smartYield.GET("/rewards/pools/:poolAddress/transactions", a.handleRewardPoolsStakingActions)

	notifs := a.engine.Group("/api/notifications")
	notifs.GET("/list", a.handleNotifications)
}
