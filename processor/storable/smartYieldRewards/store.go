package smartYieldRewards

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/barnbridge/barnbridge-backend/state"
)

func (s *Storable) storeProcessed(tx *sql.Tx) error {
	err := s.storeClaimEvents(tx)
	if err != nil {
		return err
	}

	err = s.storeStakingEvents(tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeClaimEvents(tx *sql.Tx) error {
	if len(s.processed.claims) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_rewards_claims", "user_address", "amount", "pool_address", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.claims {
		_, err = stmt.Exec(a.UserAddress, a.Amount.String(), a.PoolAddress, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

func (s *Storable) storeStakingEvents(tx *sql.Tx) error {
	if len(s.processed.stakingActions) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_rewards_staking_actions", "user_address", "amount", "balance_after", "action_type", "pool_address", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.stakingActions {
		_, err = stmt.Exec(a.UserAddress, a.Amount.String(), a.BalanceAfter.String(), a.ActionType, a.PoolAddress, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

	stmt, err = tx.Prepare(pq.CopyIn("smart_yield_transaction_history", "protocol_id", "sy_address", "underlying_token_address", "user_address", "amount", "tranche", "transaction_type", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.stakingActions {
		rewardPool := state.RewardPoolByAddress(a.PoolAddress)
		syPool := state.PoolBySmartYieldAddress(rewardPool.PoolTokenAddress)

		_, err = stmt.Exec(syPool.ProtocolId, syPool.SmartYieldAddress, syPool.UnderlyingAddress, a.UserAddress, a.Amount.String(), "JUNIOR", a.ActionType, a.TransactionHash,
			a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)

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
