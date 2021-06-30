package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddReceiverColumnToMintedNFTsTable, downAddReceiverColumnToMintedNFTsTable)
}

func upAddReceiverColumnToMintedNFTsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	alter table minted_nfts 
			add receiver text
	`)
	return err
}

func downAddReceiverColumnToMintedNFTsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table minted_nfts 
		    drop column if exists receiver
		    ;`)
	return err
}
