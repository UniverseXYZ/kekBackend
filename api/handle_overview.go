package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Overview struct {
	AvgLockTimeSeconds  int64  `json:"avgLockTimeSeconds"`
	Holders             int64  `json:"holders"`
	TotalDelegatedPower string `json:"totalDelegatedPower"`
	TotalVBond          string `json:"totalVbond"`
	Voters              int64  `json:"voters"`
	BarnUsers           int64  `json:"barnUsers"`
}

func (a *API) BondOverview(c *gin.Context) {
	var overview Overview

	err := a.db.QueryRow(`select coalesce(avg(locked_until - locked_at),0)::bigint from barn_locks;`).Scan(&overview.AvgLockTimeSeconds)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select coalesce(sum(voting_power(user_address)),0) as total_voting_power from barn_users;`).Scan(&overview.TotalVBond)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select count(*) from bond_users_with_balance where balance > 0;`).Scan(&overview.Holders)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select coalesce(sum(total),0) from ( select case when action_type = 'INCREASE' then sum(amount)
		when action_type = 'DECREASE' then -sum(amount) end total
		from barn_delegate_changes
		group by action_type ) x;
    `).Scan(&overview.TotalDelegatedPower)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select count(*)
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

	err = a.db.QueryRow(`select count(*) from barn_users;`).Scan(&overview.BarnUsers)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	OK(c, overview)
}
