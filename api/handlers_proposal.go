package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	types2 "github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) proposalDetailsHandler(c *gin.Context) {
	proposalID := utils.CleanUpHex(c.Param("proposalID"))
	var (
		Id           uint64
		Proposer     string
		Description  string
		Title        string
		CreateTime   uint64
		StartTime    uint64
		Quorum       string
		Eta          uint64
		ForVotes     string
		AgainstVotes string
		Canceled     bool
		Executed     bool
		Targets      types2.JSONStringArray
		Values       types2.JSONStringArray
		Signatures   types2.JSONStringArray
		Calldatas    types2.JSONStringArray
	)
	err := a.core.DB().QueryRow(`select proposal_ID,proposer,description,title,create_time,start_time,quorum,eta,for_votes,against_votes,canceled,executed,targets,"values",signatures,calldatas,"timestamp"
     , from governance_proposals where proposal_ID = $1 limit 1`, proposalID).Scan(&Id, &Proposer, &Description, &Title, &CreateTime, &StartTime, &Quorum, &Eta, &ForVotes, &AgainstVotes, &Canceled, &Executed, &Targets,
		&Values, &Signatures, &Calldatas)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	proposal := types.Proposal{
		Id:           Id,
		Proposer:     Proposer,
		Description:  Description,
		Title:        Title,
		CreateTime:   CreateTime,
		StartTime:    StartTime,
		Quorum:       Quorum,
		Eta:          Eta,
		ForVotes:     ForVotes,
		AgainstVotes: AgainstVotes,
		Canceled:     Canceled,
		Executed:     Executed,
		Targets:      Targets,
		Values:       Values,
		Signatures:   Signatures,
		Calldatas:    Calldatas,
	}

	OK(c, proposal)
}
