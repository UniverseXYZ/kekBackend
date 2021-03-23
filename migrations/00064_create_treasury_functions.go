package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTreasuryFunctions, downCreateTreasuryFunctions)
}

func upCreateTreasuryFunctions(tx *sql.Tx) error {
	_, err := tx.Exec(`
			create or replace function get_account_label(addr text) returns text
				language plpgsql as
			$$
			declare
				labelText text;
			begin
				select into labelText (select label from labels where address = addr);
				if labelText IS NULL then
					return '';
				end if;
				return labelText;
			end;
			$$;
			
			create or replace function get_treasury_tokens(addr text) returns table (token_address text,symbol text, decimals bigint)
				language plpgsql
			as
			$$
			declare
				sumIn numeric(78);
				sumOut numeric(78);
			begin
				 sumIn:= coalesce(sum(amount),0) from account_erc20_transfers where account = addr and tx_direction = 'IN';
   				 sumOut:= coalesce(sum(amount),0)  from account_erc20_transfers where account = addr and tx_direction = 'OUT';
				if sumIn - sumOut > 0 then
					return query
						select distinct t.token_address,
										e20t.symbol ,
										e20t.decimals
						from account_erc20_transfers t
								 inner join erc20_tokens e20t on t.token_address = e20t.token_address
						where t.account = addr;
				end if;
			end;
			$$;
	`)

	return err

}

func downCreateTreasuryFunctions(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop function if exists get_account_label;
		drop function if exists get_treasury_tokens;
	`)
	return err
}
