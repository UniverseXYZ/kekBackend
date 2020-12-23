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

	rows, err := a.core.DB().Query(`select user_address,balance_of(user_address) as bond_staked,
		   ( select coalesce(locked_until,0)
			 from barn_locks
			 where user_address = barn_users.user_address
			 order by included_in_block desc, log_index desc )                                           as locked_until,
		   delegated_power(user_address),
		   ( select count(*) from governance_votes where user_id = barn_users.user_address ) +
		   ( select count(*) from governance_cancellation_votes where user_id = barn_users.user_address ) as votes,
		   ( select count(*) from governance_proposals where proposer = barn_users.user_address )         as proposals,
		   voting_power(user_address)                                                                     as voting_power,
		   has_active_delegation(user_address)
	from barn_users
	order by voting_power desc;`)

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
	OK(c, votersList)
}
