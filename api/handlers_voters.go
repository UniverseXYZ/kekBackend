package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) VoteDetailsHandler(c *gin.Context) {
	proposalID := utils.CleanUpHex(c.Param("proposalID"))

	rows, err := a.core.DB().Query(`select proposal_ID,user_ID,support,power,"timestamp",tx_hash,tx_index,log_index,logged_by from governance_votes where proposal_ID =$1 order by "timestamp" desc`, proposalID)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	defer rows.Close()

	var votesList []types.Vote

	for rows.Next() {
		var (
			LoggedBy         string
			TransactionHash  string
			TransactionIndex int64
			LogIndex         int64

			ProposalID uint64
			User       string
			Support    bool
			Power      int64
			Timestamp  int64
		)
		err := rows.Scan(&ProposalID, &User, &Support, &Power, &Timestamp, &TransactionHash, &TransactionIndex, &LogIndex, &LoggedBy)
		if err != nil {
			Error(c, err)
			return
		}

		vote := types.Vote{
			ProposalID:       ProposalID,
			User:             User,
			Support:          Support,
			Power:            Power,
			Timestamp:        Timestamp,
			TransactionHash:  TransactionHash,
			TransactionIndex: TransactionIndex,
			LoggedBy:         LoggedBy,
			LogIndex:         LogIndex,
		}
		votesList = append(votesList, vote)
	}

	if len(votesList) == 0 {
		NotFound(c)
		return
	}

	OK(c, votesList)
}

func (a *API) VoteCanceledDetailsHandler(c *gin.Context) {
	proposalID := utils.CleanUpHex(c.Param("proposalID"))

	rows, err := a.core.DB().Query(`select proposal_ID,user_ID,"timestamp",tx_hash,tx_index,log_index,logged_by from governance_votes where proposal_ID =$1 order by "timestamp" desc`, proposalID)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	defer rows.Close()
}
