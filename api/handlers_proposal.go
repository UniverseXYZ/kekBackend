package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	types2 "github.com/barnbridge/barnbridge-backend/types"
)

func (a *API) ProposalDetailsHandler(c *gin.Context) {
	proposalID := c.Param("proposalID")

	var (
		id          uint64
		proposer    string
		description string
		title       string
		createTime  uint64
		targets     types2.JSONStringArray
		values      types2.JSONStringArray
		signatures  types2.JSONStringArray
		calldatas   types2.JSONStringArray
		timestamp   int64
	)
	err := a.core.DB().QueryRow(`select proposal_id,proposer,description,title,create_time,targets,"values",signatures,calldatas,block_timestamp from governance_proposals where proposal_ID = $1 `, proposalID).Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	proposal := types.Proposal{
		Id:             id,
		Proposer:       proposer,
		Description:    description,
		Title:          title,
		CreateTime:     createTime,
		Targets:        targets,
		Values:         values,
		Signatures:     signatures,
		Calldatas:      calldatas,
		BlockTimestamp: timestamp,
	}

	OK(c, proposal)
}

func (a *API) AllProposalHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "10")

	rows, err := a.core.DB().Query(`select proposal_ID,proposer,description,title,create_time,targets,"values",signatures,calldatas,block_timestamp from governance_proposals where proposal_id <= $1 order by block_timestamp desc limit $2`, offset, limit)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var proposalList []types.Proposal

	for rows.Next() {
		var (
			id          uint64
			proposer    string
			description string
			title       string
			createTime  uint64
			targets     types2.JSONStringArray
			values      types2.JSONStringArray
			signatures  types2.JSONStringArray
			calldatas   types2.JSONStringArray
			timestamp   int64
		)
		err := rows.Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp)
		if err != nil {
			Error(c, err)
			return
		}

		proposal := types.Proposal{
			Id:             id,
			Proposer:       proposer,
			Description:    description,
			Title:          title,
			CreateTime:     createTime,
			Targets:        targets,
			Values:         values,
			Signatures:     signatures,
			Calldatas:      calldatas,
			BlockTimestamp: timestamp,
		}
		proposalList = append(proposalList, proposal)
	}

	if len(proposalList) == 0 {
		NotFound(c)
		return
	}

	OK(c, proposalList)
}
