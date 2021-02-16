package smartYield

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (s *Storable) storeProcessed(tx *sql.Tx) error {
	err := s.storeTokenBuyTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store token buy trades")
	}

	err = s.storeTokenSellTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store token sell trades")
	}

	err = s.storeJuniorBuyTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store junior buy trades")
	}

	err = s.storeJuniorRedeemTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store junior redeem trades")
	}

	err = s.storeSeniorBuyTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store senior buy trades")
	}

	err = s.storeSeniorRedeemTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store senior redeem trades")
	}

	return nil
}

func (s *Storable) storeTokenBuyTrades(tx *sql.Tx) error {
	if len(s.processed.tokenActions.tokenBuyTrades) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_token_buy", "buyer_address", "underlying_in", "tokens_out", "fee", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.tokenActions.tokenBuyTrades {
		_, err = stmt.Exec(a.BuyerAddress, a.UnderlyingIn.String(), a.TokensOut.String(), a.Fee.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeTokenSellTrades(tx *sql.Tx) error {
	if len(s.processed.tokenActions.tokenSellTrades) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_token_sell", "seller_address", "tokens_in", "underlying_out", "forfeits", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.tokenActions.tokenSellTrades {
		_, err = stmt.Exec(a.SellerAddress, a.TokensIn.String(), a.UnderlyingOut.String(), a.Forfeits.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeJuniorBuyTrades(tx *sql.Tx) error {
	if len(s.processed.juniorActions.juniorBondBuys) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_junior_buy", "buyer_address", "junior_bond_id", "tokens_in", "matures_at", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.juniorActions.juniorBondBuys {
		_, err = stmt.Exec(a.BuyerAddress, a.JuniorBondID.String(), a.TokensIn.String(), a.MaturesAt.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeJuniorRedeemTrades(tx *sql.Tx) error {
	if len(s.processed.juniorActions.juniorBondRedeems) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_junior_redeem", "owner_address", "junior_bond_id", "underlying_out", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.juniorActions.juniorBondRedeems {
		_, err = stmt.Exec(a.OwnerAddress, a.JuniorBondID.String(), a.UnderlyingOut.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeSeniorBuyTrades(tx *sql.Tx) error {
	if len(s.processed.seniorActions.seniorBondBuys) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_senior_buy", "buyer_address", "senior_bond_id", "underlying_in", "gain", "for_days", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.seniorActions.seniorBondBuys {
		_, err = stmt.Exec(a.BuyerAddress, a.SeniorBondID.String(), a.UnderlyingIn.String(), a.Gain.String(), a.ForDays.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeSeniorRedeemTrades(tx *sql.Tx) error {
	if len(s.processed.seniorActions.seniorBondRedeems) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_senior_buy", "owner_address", "senior_bond_id", "fee", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.seniorActions.seniorBondRedeems {
		_, err = stmt.Exec(a.OwnerAddress, a.SeniorBondID.String(), a.Fee.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeJTokenTransfers(tx *sql.Tx) error {
	if len(s.processed.jTokenTransfers) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("jtoken_transfers", "tx_hash", "tx_index", "log_index", "sender", "receiver", "value", "included_in_block", "block_timestamp"))
	if err != nil {
		return err
	}

	for _, t := range s.processed.jTokenTransfers {
		_, err = stmt.Exec(t.TransactionHash, t.TransactionIndex, t.LogIndex, t.From, t.To, t.Value.String(), s.processed.blockNumber, s.processed.blockTimestamp)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeSmartBondTransfers(tx *sql.Tx) error {
	if len(s.processed.smartBondTransfers) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("smart_bond_transfers", "tx_hash", "tx_index", "log_index", "token_address", "sender", "receiver", "token_id", "included_in_block", "block_timestamp"))
	if err != nil {
		return err
	}

	for _, t := range s.processed.smartBondTransfers {
		_, err = stmt.Exec(t.TransactionHash, t.TransactionIndex, t.LogIndex, t.TokenAddress, t.Sender, t.Receiver, t.TokenID.String(), s.processed.blockNumber, s.processed.blockTimestamp)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}
