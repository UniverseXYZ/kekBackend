package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableGovernanceAbrogationProposalsVotes, downCreateTableGovernanceAbrogationProposalsVotes)
}

func upCreateTableGovernanceAbrogationProposalsVotes(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table governance_abrogation_votes
	(
		proposal_id bigint not null,
		user_id text not null,
		support boolean not null,
		power numeric(78) not null,
		block_timestamp bigint,
		tx_hash text not null,
		tx_index integer not null,
		log_index integer not null,
		logged_by text not null,
		included_in_block bigint not null,
		created_at timestamp default now()
	);

	create index governance_abrogation_votes_proposal_id_idx
		on governance_abrogation_votes (proposal_id desc);

	create index governance_abrogation_votes_proposal_id_composed_idx
		on governance_abrogation_votes (proposal_id asc, user_id asc, block_timestamp desc);

	create index governance_abrogation_votes_user_id_idx
		on governance_abrogation_votes (lower(user_id));
	`)
	return err
}

func downCreateTableGovernanceAbrogationProposalsVotes(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists governance_abrogation_votes")
	return err
}
