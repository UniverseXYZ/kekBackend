package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddReceiverColumnToUniverseTable, downAddReceiverColumnToUniverseTable)
}

func upAddReceiverColumnToUniverseTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	alter table universe 
			add owner text
	`)
	return err
}

func downAddReceiverColumnToUniverseTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table universe 
		    drop column if exists owner
		    ;`)
	return err
}
