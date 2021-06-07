package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableAuctions, downCreateTableAuctions)
}

func upCreateTableAuctions(tx *sql.Tx) error {
	_, err := tx.Exec(`

	create table auctions
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		data JSONB not null,
		block_timestamp bigint not null,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	`)
	return err
}

func downCreateTableAuctions(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists auctions")
	return err
}
