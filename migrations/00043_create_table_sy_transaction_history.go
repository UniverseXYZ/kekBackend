package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSyTransactionHistory, downCreateTableSyTransactionHistory)
}

func upCreateTableSyTransactionHistory(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create type tx_type as enum('JUNIOR_DEPOSIT','JUNIOR_INSTANT_WITHDRAW','JUNIOR_REGULAR_WITHDRAW','JUNIOR_REDEEM','SENIOR_DEPOSIT','SENIOR_REDEEM','JTOKEN_TRANSFER','JBOND_TRANSFER','SBOND_TRANSFER');
		create table smart_yield_transaction_history
		(
		    protocol_id       				text    not null,
			sy_address        				text    not null,
			underlying_token_address 		text not null ,
			user_address 					text not null ,
			amount 							bigint,
			tranche 						text not null,
			transaction_type 				tx_type not null,
			
		
			tx_hash           text    not null,
			tx_index          integer not null,
			log_index         integer not null,
		
			block_timestamp   bigint  not null,
			included_in_block bigint  not null,
		
			created_at        timestamp default now()
		);
	`)

	return err
}

func downCreateTableSyTransactionHistory(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists smart_yield_transaction_history")
	return err
}
