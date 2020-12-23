package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Overview struct {
	AvgLockTimeSeconds  string
	Holders             string
	TotalDelegatedPower string
	TotalVBond          string
	Voters              string
	BarnUsers           string
}

func (a *API) BondOverview(c *gin.Context) {
	var overview Overview

	err := a.core.DB().QueryRow(`select avg(locked_until - locked_at) from barn_locks;`).Scan(&overview.AvgLockTimeSeconds)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.core.DB().QueryRow(`select sum(voting_power(user_address)) as total_voting_power from barn_users;`).Scan(&overview.TotalVBond)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.core.DB().QueryRow(`select count(*) from bond_users_with_balance where balance > 0;`).Scan(&overview.Holders)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.core.DB().QueryRow(`select sum(total) from ( select case when action_type = 'INCREASE' then sum(amount)
		when action_type = 'DECREASE' then -sum(amount) end total
		from barn_delegate_changes
		group by action_type ) x;
    `).Scan(&overview.TotalDelegatedPower)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.core.DB().QueryRow(`select count(*)
	from ( select distinct user_id as address
		   from governance_votes
		   union
		   select distinct user_id
		   from governance_cancellation_votes ) x; 
		   `).Scan(&overview.Voters)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.core.DB().QueryRow(`select count(*) from barn_users;`).Scan(&overview.BarnUsers)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	OK(c, overview)
}
