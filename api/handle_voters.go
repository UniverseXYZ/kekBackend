package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Voter struct {
	Address             string `json:"address"`
	BondStaked          string `json:"bondStaked"`
	LockedUntil         int64  `json:"lockedUntil"`
	DelegatedPower      string `json:"delegatedPower"`
	Votes               int64  `json:"votes"`
	Proposals           int64  `json:"proposals"`
	VotingPower         string `json:"votingPower"`
	HasActiveDelegation bool   `json:"hasActiveDelegation"`
}

func (a *API) handleVoters(c *gin.Context) {
	var votersList []Voter
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
	}

	rows, err := a.core.DB().Query(` select * from voters order by voting_power desc offset $1 limit $2 ;`, offset, limit)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	for rows.Next() {
		var voter Voter
		err := rows.Scan(&voter.Address, &voter.BondStaked, &voter.LockedUntil, &voter.DelegatedPower, &voter.Votes, &voter.Proposals, &voter.VotingPower, &voter.HasActiveDelegation)
		if err != nil {
			Error(c, err)
			return
		}
		votersList = append(votersList, voter)
	}

	if len(votersList) == 0 {
		NotFound(c)
		return
	}

	var count int
	err = a.core.DB().QueryRow(`select count(*) from voters`).Scan(&count)

	OK(c, votersList, map[string]interface{}{"count": count})
}
