package migrations

import (
"database/sql"

"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddProcessedColumnToMintedNFTsTable, downAddProcessedColumnToMintedNFTsTable)
}

func upAddProcessedColumnToMintedNFTsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	alter table minted_nfts
			add id SERIAL PRIMARY KEY,
			add processed boolean DEFAULT false
	`)
	return err
}

func downAddProcessedColumnToMintedNFTsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table minted_nfts 
		    drop column if exists processed
		    ;`)
	return err
}
