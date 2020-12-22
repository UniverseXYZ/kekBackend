package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

var votesList []types.Vote
var canceledVotesList []types.VoteCanceled

func (a *API) VoteDetailsHandler(proposalID string) {

	rows, err := a.core.DB().Query(`select proposal_ID,user_ID,support,power,block_timestamp,tx_hash,tx_index,log_index,logged_by from governance_votes where proposal_ID =$1 order by "timestamp" desc`, proposalID)
	if err != nil && err != sql.ErrNoRows {
		//Error(c, err)
		return
	}
	defer rows.Close()

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
			//Error(c, err)
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
			Canceled:         false,
		}
		votesList = append(votesList, vote)
	}
	/*
		if len(votesList) == 0 {
			NotFound(c)
			return
		}

		OK(c, votesList)*/
}

func (a *API) VoteCanceledDetailsHandler(proposalID string) {

	rows, err := a.core.DB().Query(`select proposal_ID,user_ID,block_timestamp,tx_hash,tx_index,log_index,logged_by from governance_votes where proposal_ID =$1 order by "timestamp" desc`, proposalID)
	if err != nil && err != sql.ErrNoRows {
		//Error(c, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			LoggedBy         string
			TransactionHash  string
			TransactionIndex int64
			LogIndex         int64

			ProposalID uint64
			User       string
			Timestamp  int64
		)

		err := rows.Scan(&ProposalID, &User, &Timestamp, &TransactionHash, &TransactionIndex, &LogIndex, &LoggedBy)
		if err != nil {
			//Error(c, err)
			return
		}
		canceledVote := types.VoteCanceled{
			ProposalID: ProposalID,
			User:       User,
			Timestamp:  Timestamp,

			LoggedBy:         LoggedBy,
			TransactionIndex: TransactionIndex,
			TransactionHash:  TransactionHash,
			LogIndex:         LogIndex,
		}
		canceledVotesList = append(canceledVotesList, canceledVote)
	}

	/*if len(canceledVotesList) == 0 {
		NotFound(c)
		return
	}

	OK(c, canceledVotesList)*/
}

func (a *API) VotesHandler(c *gin.Context) {
	proposalID := utils.CleanUpHex(c.Param("proposalID"))

	a.VoteDetailsHandler(proposalID)
	if len(votesList) == 0 {
		NotFound(c)
		return
	}

	if len(canceledVotesList) == 0 {
		OK(c, votesList)
		return
	}

	for _, vote := range votesList {
		canceled := false
		for _, canceledVote := range canceledVotesList {
			if vote.User == canceledVote.User {
				if vote.Timestamp < canceledVote.Timestamp {
					canceled = true
				}
			}
		}
		vote.Canceled = canceled
	}

}
