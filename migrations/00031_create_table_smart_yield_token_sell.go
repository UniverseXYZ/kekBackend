package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSmartYieldTokenSell, downCreateTableSmartYieldTokenSell)
}

func upCreateTableSmartYieldTokenSell(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table smart_yield_token_sell
		(
			sy_address        text    not null,
			seller_address    text    not null,
			tokens_in         bigint,
			underlying_out    bigint,
			forfeits          bigint,
		
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

func downCreateTableSmartYieldTokenSell(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists smart_yield_token_sell")
	return err
}
