package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableBarnStakingActions, downCreateTableBarnStakingActions)
}

func upCreateTableBarnStakingActions(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table supernova_staking_actions
	(
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		address text not null,
		user_address text not null,
		action_type action_type not null,
		amount numeric(78) not null,
		balance_after numeric(78) not null,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	create index user_balance_idx
		on supernova_staking_actions (user_address asc, included_in_block desc, log_index desc);

	create index supernova_staking_actions_included_in_block_idx
		on supernova_staking_actions (included_in_block desc);
	
		create or replace function refresh_supernova_users() returns TRIGGER
			language plpgsql as
		$$
		begin
			refresh materialized view concurrently supernova_users;
			return null;
		end
		$$;

	create trigger refresh_supernova_users
		after insert or update or delete or truncate
		on supernova_staking_actions
		execute procedure refresh_supernova_users();

	`)
	return err
}

func downCreateTableBarnStakingActions(tx *sql.Tx) error {
	_, err := tx.Exec(`
        drop function refresh_supernova_users;
		drop trigger refresh_supernova_users on supernova_staking_actions;
		drop table supernova_staking_actions;
	`)
	return err
}
