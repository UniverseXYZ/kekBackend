package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableUniverse, downCreateTableUniverse)
}

func upCreateTableUniverse(tx *sql.Tx) error {
	_, err := tx.Exec(`

	create table universe
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		token_name text not null,
		token_symbol text not null,
		contract_address text not null,
		block_timestamp bigint not null,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	`)
	return err
}

func downCreateTableUniverse(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists universe")
	return err
}
