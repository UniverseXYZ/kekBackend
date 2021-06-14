package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableMonitoredNFTs, downCreateTableMonitoredNFTs)
}

func upCreateTableMonitoredNFTs(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table monitored_nfts
	(
		address text not null,
		created_at timestamp default now()
	);

	`)
	return err
}

func downCreateTableMonitoredNFTs(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists monitored_nfts;
		`)
	return err
}
