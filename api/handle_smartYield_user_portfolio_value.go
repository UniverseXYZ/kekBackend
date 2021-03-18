package api

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type SYPortfolioValuePoint struct {
	Timestamp   time.Time `json:"timestamp"`
	SeniorValue *float64  `json:"seniorValue,omitempty"`
	JuniorValue *float64  `json:"juniorValue,omitempty"`
}

func (a *API) handleSYUserPortfolioValue(c *gin.Context) {
	user, err := getQueryAddress(c, "address")
	if err != nil {
		BadRequest(c, err)
		return
	}

	rows, err := a.db.Query(`
		select ts,
			   junior_portfolio_value_at_ts($1, ts),
			   senior_portfolio_value_at_ts($1, ts)
		from generate_series(( select extract(epoch from now() - interval '30 days')::bigint ),
							 ( select extract(epoch from now()) )::bigint, 12 * 60 * 60) as ts
	`, user)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var points []SYPortfolioValuePoint

	for rows.Next() {
		var ts int64
		var junior, senior float64

		err := rows.Scan(&ts, &junior, &senior)
		if err != nil {
			Error(c, err)
			return
		}

		points = append(points, SYPortfolioValuePoint{
			Timestamp:   time.Unix(ts, 0),
			SeniorValue: &senior,
			JuniorValue: &junior,
		})
	}

	OK(c, points)
}

func (a *API) handleSYUserSeniorPortfolioValue(c *gin.Context) {
	user, err := getQueryAddress(c, "address")
	if err != nil {
		BadRequest(c, err)
		return
	}

	rows, err := a.db.Query(`
		select ts,
			   senior_portfolio_value_at_ts($1, ts)
		from generate_series(( select extract(epoch from now() - interval '30 days')::bigint ),
							 ( select extract(epoch from now()) )::bigint, 12 * 60 * 60) as ts
	`, user)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var points []SYPortfolioValuePoint

	for rows.Next() {
		var ts int64
		var senior float64

		err := rows.Scan(&ts, &senior)
		if err != nil {
			Error(c, err)
			return
		}

		points = append(points, SYPortfolioValuePoint{
			Timestamp:   time.Unix(ts, 0),
			SeniorValue: &senior,
		})
	}

	OK(c, points)
}

func (a *API) handleSYUserJuniorPortfolioValue(c *gin.Context) {
	user, err := getQueryAddress(c, "address")
	if err != nil {
		BadRequest(c, err)
		return
	}

	rows, err := a.db.Query(`
		select ts,
			   junior_portfolio_value_at_ts($1, ts)
		from generate_series(( select extract(epoch from now() - interval '30 days')::bigint ),
							 ( select extract(epoch from now()) )::bigint, 12 * 60 * 60) as ts
	`, user)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var points []SYPortfolioValuePoint

	for rows.Next() {
		var ts int64
		var junior float64

		err := rows.Scan(&ts, &junior)
		if err != nil {
			Error(c, err)
			return
		}

		points = append(points, SYPortfolioValuePoint{
			Timestamp:   time.Unix(ts, 0),
			JuniorValue: &junior,
		})
	}

	OK(c, points)
}
