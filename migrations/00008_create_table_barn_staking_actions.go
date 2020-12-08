package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableBarnEvents, downCreateTableBarnEvents)
}

func upCreateTableBarnEvents(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table bond_transfers
	(
		tx_hash                    text    not null,
		tx_index 				   integer not null,
		log_index                  integer not null,
		address					   text not null,
		userAddress				text not null,
		
		value 					   numeric (78),
		included_in_block          bigint  not null,
		created_at                 timestamp default now()
	);
	
	`)
	return err
}

func downCreateTableBarnEvents(tx *sql.Tx) error {
	_, err := tx.Exec("drop table bond_transfers")
	return err
}
