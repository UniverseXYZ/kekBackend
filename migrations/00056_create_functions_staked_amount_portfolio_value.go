package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateFunctionsStakedAmountPortfolioValue, downCreateFunctionsStakedAmountPortfolioValue)
}

func upCreateFunctionsStakedAmountPortfolioValue(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create or replace function staked_amount_at_ts(pool text ,address text, ts bigint) returns numeric(78)
			language plpgsql as
		$$
		declare
			value numeric(78);
		begin
			select into value sum(balance_after)
			from smart_yield_rewards_staking_action as a
			where a.pool_address = pool and a.block_timestamp <=ts and a.user_address = address;
			
			return value;
		end;
		$$;

	create or replace function staked_value_at_ts(pool text ,address text, ts bigint) returns  double precision
			language plpgsql as
		$$
		declare
		value double precision;
	begin
		select into value coalesce(
		    staked_amount_at_ts(pool,address,ts)::numeric(78, 18) / 
		    	pow (10,(select underlying_decimals from  smart_yield_pools as p 
		    		where p.sy_address = (select pool_token_address from smart_yield_reward_pools as r where r.pool_address = pool )  ))*
				token_underlying_price_at_ts((select pool_token_address from smart_yield_reward_pools as r where r.pool_address = pool), ts)
	
		,0);
		return value;
	end;
	$$;

	create or replace function token_underlying_price_at_ts(addr text, ts bigint) returns double precision
				language plpgsql as
			$$
			declare
				price double precision;
			begin
				select into price price_usd
				from smart_yield_prices p
				where p.protocol_id = ( select protocol_id from smart_yield_pools where sy_address = addr )
				  and p.token_address = ( select underlying_address from smart_yield_pools where sy_address = addr )
				  and block_timestamp <= ts
				order by block_timestamp desc
				limit 1;
			
				return price;
			end;
			$$;

`)
	return err
}

func downCreateFunctionsStakedAmountPortfolioValue(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop function if exists staked_amount_at_ts;
		drop function if exists token_underlying_price_at_ts;
		drop function if exists staked_value_at_ts;
	`)
	return err
}
