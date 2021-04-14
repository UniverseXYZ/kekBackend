package core

import (
	"database/sql"

	"github.com/kekDAO/kekBackend/metrics"
)

func (c *Core) DB() *sql.DB {
	return c.db
}

func (c *Core) Metrics() *metrics.Provider {
	return c.metrics
}
