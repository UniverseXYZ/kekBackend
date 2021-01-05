package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) CancellationProposalDetailsHandler(c *gin.Context) {
	proposalID := utils.CleanUpHex(c.Param("proposalID"))
	var (
		ProposalID uint64
		Creator    string
		CreateTime uint64
	)

	err := a.core.DB().QueryRow(`select proposal_id, creator ,create_time from governance_cancellation_proposals where proposal_id = $1`, proposalID).Scan(&ProposalID, &Creator, &CreateTime)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	cancellationProposal := types.CancellationProposal{
		ProposalID: ProposalID,
		Creator:    Creator,
		CreateTime: CreateTime,
	}

	OK(c, cancellationProposal)
}

func (a *API) AllCancellationProposals(c *gin.Context) {
	rows, err := a.core.DB().Query(`select proposal_id, creator ,create_time from governance_cancellation_proposals order by create_time desc `)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var cancellationProposalList []types.CancellationProposal

	for rows.Next() {
		var (
			ProposalID uint64
			Creator    string
			CreateTime uint64
		)

		err := rows.Scan(&ProposalID, &Creator, &CreateTime)
		if err != nil {
			Error(c, err)
			return
		}

		cancellationProposal := types.CancellationProposal{
			ProposalID: ProposalID,
			Creator:    Creator,
			CreateTime: CreateTime,
		}

		cancellationProposalList = append(cancellationProposalList, cancellationProposal)
	}

	if len(cancellationProposalList) == 0 {
		NotFound(c)
		return
	}

	OK(c, cancellationProposalList)
}
