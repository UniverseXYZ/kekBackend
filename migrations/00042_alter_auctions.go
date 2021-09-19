package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddProcessedColumnToAuctionsTable, downAddProcessedColumnToAuctionsTable)
}

func upAddProcessedColumnToAuctionsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table auctions
				add id SERIAL PRIMARY KEY,
				add processed boolean DEFAULT false
	`)
	return err
}

func downAddProcessedColumnToAuctionsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`alter table auctions drop column if exists processed;`)
	return err
}
