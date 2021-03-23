package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterColumnsTableAccountErc20Transfers, downAlterColumnsTableAccountErc20Transfers)
}

func upAlterColumnsTableAccountErc20Transfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create type transfer_type as enum('IN','OUT');
		alter table account_erc20_transfers rename column sender to account;
		alter table account_erc20_transfers rename column receiver to counterparty;
		alter table account_erc20_transfers rename column "value" to amount;

		alter table account_erc20_transfers add column tx_direction transfer_type;

		create or replace function get_account_label(addr text) returns text
			language plpgsql as
			$$
			declare
				labelText text;
			begin
				select into labelText (select label from labels where address = addr);
				if labelText IS NULL then
					return 'unknown';
				end if;
				return labelText;
			end;
			$$;
`)
	return err
}

func downAlterColumnsTableAccountErc20Transfers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		alter table account_erc20_transfers rename column account to sender;
		alter table account_erc20_transfers rename column counterparty to  receiver;
		alter table account_erc20_transfers rename column amount  to "value";
		alter table account_erc20_transfers  drop column if exists tx_direction;
		drop function if exists get_account_label;
	`)
	return err
}
