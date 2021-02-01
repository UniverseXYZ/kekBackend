package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) AbrogationProposalDetailsHandler(c *gin.Context) {
	proposalID := c.Param("proposalID")
	var (
		id           uint64
		creator      string
		createTime   uint64
		description  string
		forVotes     string
		againstVotes string
	)

	err := a.db.QueryRow(`
		select proposal_id, creator, create_time, description ,
		       coalesce(( select sum(power) from abrogation_proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select sum(power) from abrogation_proposal_votes(proposal_id) where support = false ), 0) as against_votes
		from governance_abrogation_proposals 
		where proposal_id = $1
	`, proposalID).Scan(&id, &creator, &createTime, &description, &forVotes, &againstVotes)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	abrogationProposal := types.AbrogationProposal{
		ProposalID:   id,
		Creator:      creator,
		CreateTime:   createTime,
		Description:  description,
		ForVotes:     forVotes,
		AgainstVotes: againstVotes,
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, abrogationProposal, map[string]interface{}{
		"block": block,
	})
}

func (a *API) AllAbrogationProposals(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	rows, err := a.db.Query(`select proposal_id, creator, create_time
		from governance_abrogation_proposals 
		order by proposal_id desc
		offset $1
		limit $2
	`, offset, limit)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var abrogationProposalList []types.AbrogationProposal

	for rows.Next() {
		var (
			id         uint64
			creator    string
			createTime uint64
		)

		err := rows.Scan(&id, &creator, &createTime)
		if err != nil {
			Error(c, err)
			return
		}

		abrogationProposal := types.AbrogationProposal{
			ProposalID: id,
			Creator:    creator,
			CreateTime: createTime,
		}

		abrogationProposalList = append(abrogationProposalList, abrogationProposal)
	}

	var count int
	err = a.db.QueryRow(`select count(*) from governance_abrogation_proposals`).Scan(&count)
	if err != nil {
		Error(c, err)
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, abrogationProposalList, map[string]interface{}{"count": count, "block": block})
}
