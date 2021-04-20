package supernova

import (
	"database/sql"
	"encoding/hex"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/utils"
)

func (b *SupernovaStorable) handleStakingActions(logs []web3types.Log, tx *sql.Tx) error {
	var stakingActions []StakingAction

	for _, log := range logs {
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
	if len(stakingActions) == 0 {
		log.WithField("handler", "staking actions").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("supernova_staking_actions", "tx_hash", "tx_index", "log_index", "address", "user_address", "action_type", "amount", "balance_after", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, a := range stakingActions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.LoggedBy, a.UserAddress, a.ActionType, a.Amount, a.BalanceAfter, b.Preprocessed.BlockNumber)
		if err != nil {
			return errors.Wrap(err, "could not execute statement")
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "could not close statement")
	}

	return nil
}

func (b *SupernovaStorable) decodeDepositEvent(log web3types.Log) (*StakingAction, error) {
	if !utils.LogIsEvent(log, b.supernovaAbi, DepositEvent) {
		return nil, nil
	}

	var deposit Deposit
	deposit.User = utils.Topic2Address(log.Topics[1])

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.supernovaAbi.UnpackIntoInterface(&deposit, DepositEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	baseLog, err := b.getBaseLog(log)
	if err != nil {
		return nil, err
	}

	return &StakingAction{
		BaseLog:      *baseLog,
		Amount:       deposit.Amount.String(),
		BalanceAfter: deposit.NewBalance.String(),
		UserAddress:  deposit.User,
		ActionType:   DEPOSIT,
	}, nil
}

func (b *SupernovaStorable) decodeWithdrawEvent(log web3types.Log) (*StakingAction, error) {
	if !utils.LogIsEvent(log, b.supernovaAbi, WithdrawEvent) {
		return nil, nil
	}

	var withdraw Withdraw
	withdraw.User = utils.Topic2Address(log.Topics[1])

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.supernovaAbi.UnpackIntoInterface(&withdraw, "Withdraw", data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	var stakingAction StakingAction

	stakingAction.LoggedBy = utils.Trim0x(log.Address)
	stakingAction.TransactionHash = log.TransactionHash

	stakingAction.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from supernova contract to int64")
	}

	stakingAction.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  supernova contract to int64")
	}

	stakingAction.Amount = withdraw.AmountWithdrew.String()
	stakingAction.BalanceAfter = withdraw.AmountLeft.String()
	stakingAction.UserAddress = withdraw.User
	stakingAction.ActionType = WITHDRAW

	return &stakingAction, nil
}
