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
		proposal_ID					bigint not null ,
		proposer					text not null,
		description					text not null,
		title						text not null,
		create_time					bigint not null,
		start_time					bigint not null,
		quorum						bigint not null,
		eta							bigint not null,
		for_votes					bigint not null,
		against_votes				bigint not null,
		canceled					bool not null,
		executed					bool not null,

		targets 					jsonb not null ,
		values						jsonb not null ,
		signatures					jsonb not null,
		calldatas					jsonb not null,
		
		included_in_block          bigint  not null,
		created_at                 timestamp default now()
	);
	
	`)
	return err
}

func downCreateTableGovernanceProposals(tx *sql.Tx) error {
	_, err := tx.Exec("drop table governance_proposals")
	return err
}
