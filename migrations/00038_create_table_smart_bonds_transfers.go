package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSmartBondsTransfers, downCreateTableSmartBondsTransfers)
}

func upCreateTableSmartBondsTransfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table smart_bond_transfers
	(
		tx_hash                    text    not null,
		tx_index 				   integer not null,
		log_index                  integer not null,
		token_address			   text not null,
		sender 					   text not null ,
		receiver 				   text not null,
		token_id				   bigint not null,
		included_in_block          bigint  not null,
		block_timestamp			   bigint not null,
		created_at                 timestamp default now()
	);
	
	`)
	return err
}

func downCreateTableSmartBondsTransfers(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists smart_bond_transfers")
	return err
}
