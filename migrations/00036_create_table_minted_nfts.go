package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableMintedNFTs, downCreateTableMintedNFTs)
}

func upCreateTableMintedNFTs(tx *sql.Tx) error {
	_, err := tx.Exec(`

	create table minted_nfts
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		token_id text not null,
		token_uri text not null,
		block_timestamp bigint not null,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	`)
	return err
}

func downCreateTableMintedNFTs(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists minted_nfts")
	return err
}
