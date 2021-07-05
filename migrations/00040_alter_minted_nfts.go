package migrations

import (
"database/sql"

"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterTableMintedNFTs, downAlterTableMintedNFTs)
}

func upAlterTableMintedNFTs(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table minted_nfts 
		add contract_address character varying NOT NULL;
	`)
	return err
}

func downAlterTableMintedNFTs(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table minted_nfts 
		drop column if exists contract_address;
	;`)
	return err
}
