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

	err := a.db.QueryRow(`select proposal_id, creator ,create_time from governance_cancellation_proposals where proposal_id = $1`, proposalID).Scan(&id, &creator, &createTime)

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
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	rows, err := a.db.Query(`select proposal_id, creator, create_time
		from governance_cancellation_proposals 
		order by proposal_id desc
		offset $1
		limit $2
	`, offset, limit)

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

	var count int
	err = a.db.QueryRow(`select count(*) from governance_cancellation_proposals`).Scan(&count)
	if err != nil {
		Error(c, err)
	}

	OK(c, cancellationProposalList, map[string]interface{}{"count": count})
}
