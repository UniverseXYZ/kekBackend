package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableRevenuesSlots, downCreateTableRevenuesSlots)
}

func upCreateTableRevenuesSlots(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table captured_slots
		(	id SERIAL PRIMARY KEY,
			tx_hash text not null,
			tx_index integer not null,
			log_index integer not null,
			data JSONB not null,
			block_timestamp bigint not null,
			included_in_block bigint not null,
			processed boolean DEFAULT false,
			sender text not null,
			created_at timestamp default now()
		);
	`)
	return err
}

func downCreateTableRevenuesSlots(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists captured_slots")
	return err
}
