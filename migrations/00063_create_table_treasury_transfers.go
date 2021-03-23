package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableTreasuryTransfers, downCreateTableTreasuryTransfers)
}

func upCreateTableTreasuryTransfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create type treasury_type as enum('AMOUNT_IN','AMOUNT_OUT');
		create table treasury_transfers(
		    token_address		text not null,
		    address            	text    not null,
			action_type        	treasury_type    not null,
			value             	numeric(78),
		    
		    tx_hash				text not null,
		    tx_index          	integer not null,
			log_index         	integer not null,
		
			block_timestamp   bigint  not null,
			included_in_block bigint  not null
		)
	`)
	return err
}

func downCreateTableTreasuryTransfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists treasury_transfers;
		`)
	return err
}
