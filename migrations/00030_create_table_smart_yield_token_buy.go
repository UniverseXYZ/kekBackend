package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSmartYieldTokenBuy, downCreateTableSmartYieldTokenBuy)
}

func upCreateTableSmartYieldTokenBuy(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table smart_yield_token_buy
		(
			buyer_address 		text not null,
			underlying_in 		bigint ,
			tokens_out 			bigint,
			fee 				bigint,
			
			tx_hash           	text    not null,
			tx_index 			integer not null,
			log_index           integer not null,
			block_timestamp		bigint not null,
			included_in_block	bigint  not null,
			created_at			timestamp default now()
		);
	`)

	return err

}

func downCreateTableSmartYieldTokenBuy(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists smart_yield_token_buy")
	return err
}
