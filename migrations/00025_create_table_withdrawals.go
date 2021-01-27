package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableWithdrawals, downCreateTableWithdrawals)
}

func upCreateTableWithdrawals(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table withdrawals
	(
		tx_hash                    text    not null,
		tx_index 				   integer not null,
		log_index                  integer not null,
		user_address 			   text not null ,
		token_address			   text not null,
		amount 					   numeric (78),
		included_in_block          bigint  not null,
		created_at                 timestamp default now()
	);
	
	`)
	return err
}

func downCreateTableWithdrawals(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists withdrawals")
	return err
}
