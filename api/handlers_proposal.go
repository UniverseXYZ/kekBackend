package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	types2 "github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) ProposalDetailsHandler(c *gin.Context) {
	proposalID := utils.CleanUpHex(c.Param("proposalID"))
	var (
		Id          uint64
		Proposer    string
		Description string
		Title       string
		CreateTime  uint64
		Targets     types2.JSONStringArray
		Values      types2.JSONStringArray
		Signatures  types2.JSONStringArray
		Calldatas   types2.JSONStringArray
		Timestamp   int64
	)
	err := a.core.DB().QueryRow(`select proposal_id,proposer,description,title,create_time,targets,"values",signatures,calldatas,block_timestamp from governance_proposals where proposal_ID = $1 `, proposalID).Scan(&Id, &Proposer, &Description, &Title, &CreateTime, &Targets, &Values, &Signatures, &Calldatas, &Timestamp)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	proposal := types.Proposal{
		Id:             Id,
		Proposer:       Proposer,
		Description:    Description,
		Title:          Title,
		CreateTime:     CreateTime,
		Targets:        Targets,
		Values:         Values,
		Signatures:     Signatures,
		Calldatas:      Calldatas,
		BlockTimestamp: Timestamp,
	}

	OK(c, proposal)
}

func (a *API) AllProposalHandler(c *gin.Context) {
	rows, err := a.core.DB().Query(`select proposal_ID,proposer,description,title,create_time,targets,"values",signatures,calldatas,block_timestamp from governance_proposals order by block_timestamp desc`)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var proposalList []types.Proposal

	for rows.Next() {
		var (
			Id          uint64
			Proposer    string
			Description string
			Title       string
			CreateTime  uint64
			Targets     types2.JSONStringArray
			Values      types2.JSONStringArray
			Signatures  types2.JSONStringArray
			Calldatas   types2.JSONStringArray
			Timestamp   int64
		)
		err := rows.Scan(&Id, &Proposer, &Description, &Title, &CreateTime, &Targets, &Values, &Signatures, &Calldatas, &Timestamp)
		if err != nil {
			Error(c, err)
			return
		}

		proposal := types.Proposal{
			Id:             Id,
			Proposer:       Proposer,
			Description:    Description,
			Title:          Title,
			CreateTime:     CreateTime,
			Targets:        Targets,
			Values:         Values,
			Signatures:     Signatures,
			Calldatas:      Calldatas,
			BlockTimestamp: Timestamp,
		}
		proposalList = append(proposalList, proposal)
	}

	if len(proposalList) == 0 {
		NotFound(c)
		return
	}

	OK(c, proposalList)
}
