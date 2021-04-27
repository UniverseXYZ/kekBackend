package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableLKekTransfers, downCreateTableKekTransfers)
}

func upCreateTableLKekTransfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table kek_transfers
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		sender text not null,
		receiver text not null,
		value numeric(78),
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	`)
	return err
}

func downCreateTableKekTransfers(tx *sql.Tx) error {
	_, err := tx.Exec("drop table kek_transfers")
	return err
}
