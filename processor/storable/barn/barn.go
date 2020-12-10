package barn

import (
	"database/sql"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type BarnStorable struct {
	config         Config
	Raw            *types.RawData
	barnAbi        abi.ABI
	BlockTimestamp int64
	BlockNumber    int64
}

type Withdraw struct {
	AmountWithdrew *big.Int
	AmountLeft     *big.Int
	User           string
}

type Deposit struct {
	Amount     *big.Int
	NewBalance *big.Int
	User       string
}

type StakingAction struct {
	Address          string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
	UserAddress      string
	ActionType       int
	Amount           string
	BalanceAfter     string
}

const (
	DEPOSIT = iota
	WITHDRAW
)

func NewBarnStorable(config Config, raw *types.RawData, barnAbi abi.ABI) *BarnStorable {
	return &BarnStorable{
		config:  config,
		Raw:     raw,
		barnAbi: barnAbi,
	}
}

func (b BarnStorable) ToDB(tx *sql.Tx) error {
	var stakingActions []StakingAction
	var err error

	b.BlockNumber, err = strconv.ParseInt(b.Raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "unable to process block number")
	}

	b.BlockTimestamp, err = strconv.ParseInt(b.Raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse block timestamp")
	}

	for _, data := range b.Raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(b.config.BarnAddress) {
				continue
			}

			if len(log.Topics) == 0 {
				continue
			}

			stakingActionDeposit, err := b.decodeDepositEvent(log)
			if err != nil {
				return err
			}

			if stakingActionDeposit != nil {
				stakingActions = append(stakingActions, *stakingActionDeposit)
				continue
			}

			stakingActionWithdraw, err := b.decodeWithdrawEvent(log)
			if err != nil {
				return err
			}

			if stakingActionWithdraw != nil {
				stakingActions = append(stakingActions, *stakingActionWithdraw)
				continue
			}
		}
	}
	stmt, err := tx.Prepare(pq.CopyIn("barn_staking_actions", "tx_hash", "tx_index", "log_index", "address", "user_address", "action_type", "amount", "balance_after", "included_in_block", "created_at"))
	if err != nil {
		return err
	}

	for _, a := range stakingActions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.Address, a.UserAddress, a.ActionType, a.Amount, a.BalanceAfter, b.BlockNumber, b.BlockTimestamp)
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
