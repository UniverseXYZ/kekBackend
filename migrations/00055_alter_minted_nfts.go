package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddCreatorColumnToUniverseTable, downAddCreatorColumnToUniverseTable)
}

func upAddCreatorColumnToUniverseTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	alter table minted_nfts 
			add creator text
	`)
	return err
}

func downAddCreatorColumnToUniverseTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table minted_nfts 
		    drop column if exists creator
		    ;`)
	return err
}
