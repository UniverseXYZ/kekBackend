package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddColumnsToUniverseTable, downAddColumnsToUniverseTable)
}

func upAddColumnsToUniverseTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table universe
		add id SERIAL PRIMARY KEY,
		add processed boolean DEFAULT false;
	`)
	return err
}

func downAddColumnsToUniverseTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table universe
		drop column if exists processed;
	`)
	return err
}
