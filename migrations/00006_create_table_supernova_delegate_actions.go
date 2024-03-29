package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSupernovaDelegateActions, downCreateTableSupernovaDelegateActions)
}

func upCreateTableSupernovaDelegateActions(tx *sql.Tx) error {
	_, err := tx.Exec(`

	drop type if exists delegate_action_type;
	create type delegate_action_type as enum('START','STOP');

	create table supernova_delegate_actions
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		logged_by text not null,
		sender text not null,
		receiver text not null,
		action_type delegate_action_type not null,
		timestamp bigint not null,
		included_in_block bigint not null,
		created_at timestamp default now()
	);
	
	create index user_delegation_idx
		on supernova_delegate_actions (sender asc, included_in_block desc, log_index desc);

	`)
	return err
}

func downCreateTableSupernovaDelegateActions(tx *sql.Tx) error {
	_, err := tx.Exec("drop table supernova_delegate_actions")
	return err
}
