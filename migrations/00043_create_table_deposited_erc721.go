package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableDepositedErc721, downCreateTableDepositedErc721)
}

func upCreateTableDepositedErc721(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table deposited_erc721
		(	id SERIAL PRIMARY KEY,
			tx_hash text not null,
			tx_index integer not null,
			log_index integer not null,
			data JSONB not null,
			block_timestamp bigint not null,
			included_in_block bigint not null,
			processed boolean DEFAULT false,
			created_at timestamp default now()
		);
	`)
	return err
}

func downCreateTableDepositedErc721(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists deposited_erc721")
	return err
}
