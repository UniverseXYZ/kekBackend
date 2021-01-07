package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type ProposalState struct {
	ProposalID   uint64 `json:"proposal_id"`
	EventType    string `json:"event_type"`
	Caller       string `json:"caller"`
	ForVotes     uint64 `json:"for_votes"`
	AgainstVotes uint64 `json:"against_votes"`
}

func (a *API) handleProposalState(c *gin.Context) {
	proposalID := c.Param("proposalID")

	var proposalState ProposalState
	err := a.core.DB().QueryRow(`select proposal_id ,caller,event_type from governance_events where proposal_id = $1 order by block_timestamp desc limit 1`, proposalID).Scan(&proposalState.ProposalID, &proposalState.Caller, &proposalState.EventType)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	err = a.core.DB().QueryRow(`with votes as (select distinct user_id,
		first_value(support) over (partition by user_id order by block_timestamp desc) as support,
		first_value(block_timestamp) over (partition by user_id order by block_timestamp desc) as block_timestamp,
		power
		from governance_votes
		where proposal_id = $1
		and ( select count(*)
			from governance_votes_canceled
			where governance_votes_canceled.proposal_id = governance_votes.proposal_id
			and governance_votes_canceled.user_id = governance_votes.user_id
			and governance_votes_canceled.block_timestamp > governance_votes.block_timestamp ) = 0) 
		select (select coalesce(sum(power),0) from votes where support = true) as for_votes, 
		(select coalesce(sum(power),0) from votes where support = false) as against_votes;`, proposalID).Scan(&proposalState.ForVotes, &proposalState.AgainstVotes)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	OK(c, proposalState)

}

/*func getProposalState(proposal types.Proposal) (string, error) {

	return nil, nil
}*/
