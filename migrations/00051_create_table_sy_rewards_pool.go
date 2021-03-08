package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableSyRewardsPool, downCreateTableSyRewardsPool)
}

func upCreateTableSyRewardsPool(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table smart_yield_reward_pools
		(
		    pool_address			text not null,
		    pool_token_address  	text not null,
		    reward_token_address 	text not null,
		    
			created_at            timestamp default now()
		)
	`)
	return err
}

func downCreateTableSyRewardsPool(tx *sql.Tx) error {
	_, err := tx.Exec("drop table if exists smart_yield_reward_pools")
	return err
}
