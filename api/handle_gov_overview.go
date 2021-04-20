package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Overview struct {
	AvgLockTimeSeconds     int64  `json:"avgLockTimeSeconds"`
	Holders                int64  `json:"holders"`
	TotalDelegatedPower    string `json:"totalDelegatedPower"`
	TotalVKek              string `json:"TotalVKek"`
	Voters                 int64  `json:"voters"`
	SupernovaUsers         int64  `json:"supernovaUsers"`
	HoldersStakingExcluded int64  `json:"holdersStakingExcluded"`
}

func (a *API) KekOverview(c *gin.Context) {
	var overview Overview

	err := a.db.QueryRow(`select coalesce(avg(locked_until - locked_at),0)::bigint from supernova_locks;`).Scan(&overview.AvgLockTimeSeconds)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select coalesce(sum(voting_power(user_address)),0) as total_voting_power from supernova_users;`).Scan(&overview.TotalVKek)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select count(*) from kek_users_with_balance where balance > 0;`).Scan(&overview.Holders)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select count(*) from kek_users_with_balance_no_staking where balance > 0;`).Scan(&overview.HoldersStakingExcluded)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select coalesce(sum(total),0) from ( select case when action_type = 'INCREASE' then sum(amount)
		when action_type = 'DECREASE' then -sum(amount) end total
		from supernova_delegate_changes
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
		   from governance_abrogation_votes ) x; 
		   `).Scan(&overview.Voters)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	err = a.db.QueryRow(`select count(*) from voters where kek_staked + voting_power > 0`).Scan(&overview.SupernovaUsers)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, overview, map[string]interface{}{
		"block": block,
	})
}
