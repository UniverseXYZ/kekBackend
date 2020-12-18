package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableGovernanceVotesCanceled, downCreateTableGovernanceVotesCanceled)
}

func upCreateTableGovernanceVotesCanceled(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table governance_votes_canceled
	(
		proposal_ID				   bigint not null ,
		user_ID					   text not null ,
		timestamp				   bigint,
		
		tx_hash                    text    not null,
		tx_index                   integer not null,
		log_index                  integer not null,
		logged_by                  text    not null,
		
		included_in_block          bigint  not null,
		created_at                 timestamp default now()
	);
	create unique index on governance_votes (proposal_ID,user_ID)
	`)
	return err
}

func downCreateTableGovernanceVotesCanceled(tx *sql.Tx) error {
	_, err := tx.Exec("drop table governance_votes_canceled")
	return err
}

//
