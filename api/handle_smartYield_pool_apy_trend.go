package api

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

type APYTrendPoint struct {
	Point     time.Time `json:"point"`
	SeniorAPY float64   `json:"seniorApy"`
	JuniorAPY float64   `json:"juniorApy"`
}

func (a *API) handlePoolAPYTrend(c *gin.Context) {
	pool := c.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		BadRequest(c, errors.New("invalid pool address"))
		return
	}

	rows, err := a.db.Query(`
		select date_trunc('day', block_timestamp) as scale,
			   avg(senior_apy) as senior_apy,
			   avg(junior_apy) as junior_apy
		from smart_yield_state
		where pool_address = $1
		and block_timestamp > now() - interval '7 days'
		group by scale
		order by scale;
	`, poolAddress)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var points []APYTrendPoint
	for rows.Next() {
		var p APYTrendPoint

		err := rows.Scan(&p.Point, &p.SeniorAPY, &p.JuniorAPY)
		if err != nil {
			Error(c, err)
			return
		}

		points = append(points, p)
	}

	OK(c, points)
}
