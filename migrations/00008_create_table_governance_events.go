package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableGovernanceEvents, downCreateTableGovernanceEvents)
}

func upCreateTableGovernanceEvents(tx *sql.Tx) error {
	_, err := tx.Exec(`
	drop type if exists event_type;
	create type event_type as enum('CREATED','QUEUED','EXECUTED','CANCELED');

	create table governance_events
	(
		proposal_id bigint not null,
		caller text not null,
		event_type event_type not null,
		block_timestamp bigint,
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		logged_by text not null,
		event_data jsonb,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	create index governance_votes_proposal_id_event_type_idx
		on governance_events (proposal_id, event_type);
	`)
	return err
}

func downCreateTableGovernanceEvents(tx *sql.Tx) error {
	_, err := tx.Exec("drop table governance_events")
	return err
}
