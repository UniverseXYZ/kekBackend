package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(UpBarnUsersMatView, DownBarnUsersMatView)
}

func UpBarnUsersMatView(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create materialized view supernova_users as
	select distinct user_address
	from supernova_staking_actions
	union
	select distinct receiver
	from supernova_delegate_changes;

	create unique index on supernova_users(user_address);

	`)
	return err
}

func DownBarnUsersMatView(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop materialized view supernova_users;
	`)
	return err
}
