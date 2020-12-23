package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Overview struct {
	AvgLockTimeSeconds  int64
	Holders             int64
	TotalDelegatedPower int64
	TotalVBond          int64
	Voters              int64
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

}
