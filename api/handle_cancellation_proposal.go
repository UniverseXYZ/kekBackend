package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) CancellationProposalDetailsHandler(c *gin.Context) {
	proposalID := c.Param("proposalID")
	var (
		id         uint64
		creator    string
		createTime uint64
	)

	err := a.core.DB().QueryRow(`select proposal_id, creator ,create_time from governance_cancellation_proposals where proposal_id = $1`, proposalID).Scan(&id, &creator, &createTime)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	cancellationProposal := types.CancellationProposal{
		ProposalID: id,
		Creator:    creator,
		CreateTime: createTime,
	}

	OK(c, cancellationProposal)
}

func (a *API) AllCancellationProposals(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "10")

	rows, err := a.core.DB().Query(`select proposal_id, creator ,create_time from governance_cancellation_proposals where proposal_id <= $1 order by create_time desc limit $2`, offset, limit)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var cancellationProposalList []types.CancellationProposal

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

		cancellationProposal := types.CancellationProposal{
			ProposalID: id,
			Creator:    creator,
			CreateTime: createTime,
		}

		cancellationProposalList = append(cancellationProposalList, cancellationProposal)
	}

	if len(cancellationProposalList) == 0 {
		NotFound(c)
		return
	}

	OK(c, cancellationProposalList)
}
