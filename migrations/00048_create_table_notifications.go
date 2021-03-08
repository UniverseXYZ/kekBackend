package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableNotifications, downCreateTableNotifications)
}

func upCreateTableNotifications(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table "notifications"
		(
			"target"             text,
			"type"               text      not null,
			"triggered_by_block" bigint,
			"starts_on"          timestamp,
			"expires_on"         timestamp not null,
			"message"            text,
			"metadata"           jsonb,
			"created_on"		 timestamp default now()
		)
		;
		
		create index "notifications_target_starts_on_index"
			on "notifications" ("target" asc, "starts_on" desc)
		;
		
		create index "notifications_triggered_by_block_index"
			on "notifications" ("triggered_by_block" desc)
		;
	`)

	return err
}

func downCreateTableNotifications(tx *sql.Tx) error {
	_, err := tx.Exec(`drop table if exists "notifications";`)
	return err
}
