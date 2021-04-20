package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableAccountErc20Transfers, downCreateTableAccountErc20Transfers)
}

func upCreateTableAccountErc20Transfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create type transfer_type as enum('IN','OUT');
	create table account_erc20_transfers
	(
		token_address text not null,
		account text not null,
		counterparty text not null,
		amount numeric(78),
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		block_timestamp bigint not null,
		included_in_block bigint not null,
		tx_direction transfer_type
	);

	create index account_erc20_transfers_account_addr_idx
		on account_erc20_transfers (account asc, included_in_block desc, tx_index desc, log_index desc);
	`)
	return err
}

func downCreateTableAccountErc20Transfers(tx *sql.Tx) error {
	_, err := tx.Exec(`drop table if exists account_erc20_transfers;`)
	return err
}
