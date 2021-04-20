package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableBarnDelegateChange, downCreateTableBarnDelegateChange)
}

func upCreateTableBarnDelegateChange(tx *sql.Tx) error {
	_, err := tx.Exec(`

	drop type if exists delegate_change_type;
	create type delegate_change_type as enum('INCREASE','DECREASE');

	create table supernova_delegate_changes
	(
		tx_hash                    text    not null,
		tx_index                   integer not null,
		log_index                  integer not null,
		logged_by                  text not null,
		
		action_type                   delegate_change_type not null,
		sender                        text not null,
		receiver                      text not null,
		amount                        numeric(78) not null,
		receiver_new_delegated_power  numeric(78) not null,
		timestamp                     bigint,
		
		included_in_block          bigint  not null,
		created_at                 timestamp default now()
	);

	create index user_delegated_power_idx
		on supernova_delegate_changes (receiver asc, included_in_block desc, log_index desc);

	create trigger refresh_supernova_users
		after insert or update or delete or truncate
		on supernova_delegate_changes
		execute procedure refresh_supernova_users();

	`)
	return err
}

func downCreateTableBarnDelegateChange(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop trigger refresh_supernova_users on supernova_delegate_changes;
		drop table supernova_delegate_changes;
	`)
	return err
}
