package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableRevenuesCaptured, downCreateTableRevenuesCaptured)
}

func upCreateTableRevenuesCaptured(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table captured_revenues
		(	id SERIAL PRIMARY KEY,
			tx_hash text not null,
			tx_index integer not null,
			log_index integer not null,
			data JSONB not null,
			block_timestamp bigint not null,
			included_in_block bigint not null,
			processed boolean DEFAULT false,
			created_at timestamp default now()
		);
	`)
	return err
}

func downCreateTableRevenuesCaptured(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists captured_revenues")
	return err
}
