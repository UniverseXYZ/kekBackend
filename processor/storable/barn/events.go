package barn

import (
	"encoding/hex"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

const DepositEvent = "Deposit"
const WithdrawEvent = "Withdraw"

func (b *BarnStorable) decodeDepositEvent(log web3types.Log) (*StakingAction, error) {
	if utils.CleanUpHex(log.Topics[0]) != utils.CleanUpHex(b.barnAbi.Events[DepositEvent].ID.String()) {
		return nil, nil
	}

	var deposit Deposit
	deposit.User = utils.Topic2Address(log.Topics[1])

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.barnAbi.UnpackIntoInterface(&deposit, DepositEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	var stakingAction StakingAction

	stakingAction.Address = utils.Trim0x(log.Address)
	stakingAction.TransactionHash = log.TransactionHash

	stakingAction.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from barn contract to int64")
	}

	stakingAction.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  barn contract to int64")
	}

	stakingAction.Amount = deposit.Amount.String()
	stakingAction.BalanceAfter = deposit.NewBalance.String()
	stakingAction.UserAddress = deposit.User
	stakingAction.ActionType = DEPOSIT

	return &stakingAction, nil
}

func (b *BarnStorable) decodeWithdrawEvent(log web3types.Log) (*StakingAction, error) {
	if utils.CleanUpHex(log.Topics[0]) != utils.CleanUpHex(b.barnAbi.Events[WithdrawEvent].ID.String()) {
		return nil, nil
	}

	var withdraw Withdraw
	withdraw.User = utils.Topic2Address(log.Topics[1])

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.barnAbi.UnpackIntoInterface(&withdraw, "Withdraw", data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	var stakingAction StakingAction

	stakingAction.Address = utils.Trim0x(log.Address)
	stakingAction.TransactionHash = log.TransactionHash

	stakingAction.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from barn contract to int64")
	}

	stakingAction.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  barn contract to int64")
	}

	stakingAction.Amount = withdraw.AmountWithdrew.String()
	stakingAction.BalanceAfter = withdraw.AmountLeft.String()
	stakingAction.UserAddress = withdraw.User
	stakingAction.ActionType = WITHDRAW

	return &stakingAction, nil
}
