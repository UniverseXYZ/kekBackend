package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableGovernanceProposals, downCreateTableGovernanceProposals)
}

func upCreateTableGovernanceProposals(tx *sql.Tx) error {
	_, err := tx.Exec(`

	create table governance_proposals
	(
		proposal_id bigint not null,
		proposer text not null,
		description text not null,
		title text not null,
		create_time bigint not null,
		targets jsonb not null,
		values jsonb not null,
		signatures jsonb not null,
		calldatas jsonb not null,
		block_timestamp bigint,
		included_in_block bigint not null,
		created_at timestamp default now(),
		warm_up_duration bigint,
		active_duration bigint,
		queue_duration bigint,
		grace_period_duration bigint,
		acceptance_threshold bigint,
		min_quorum bigint
	);
	
	create index governance_proposals_proposal_id_idx
		on governance_proposals (proposal_id desc);
	
	create index governance_proposals_proposer_idx
		on governance_proposals (lower(proposer));
	`)
	return err
}

func downCreateTableGovernanceProposals(tx *sql.Tx) error {
	_, err := tx.Exec("drop table governance_proposals")
	return err
}
