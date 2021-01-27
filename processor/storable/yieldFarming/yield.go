package yieldFarming

import (
	"database/sql"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type Storable struct {
	config   Config
	raw      *types.RawData
	yieldAbi abi.ABI
}

func NewStorable(config Config, raw *types.RawData, yieldAbi abi.ABI) *Storable {
	return &Storable{
		config:   config,
		raw:      raw,
		yieldAbi: yieldAbi,
	}
}

func (y *Storable) ToDB(tx *sql.Tx) error {
	var yieldDepositLogs []web3types.Log
	var yieldWithdrawLogs []web3types.Log

	for _, data := range y.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(y.config.Address) {
				continue
			}
			if len(log.Topics) == 0 {
				continue
			}

			if utils.LogIsEvent(log, y.yieldAbi, Deposit) {
				yieldDepositLogs = append(yieldDepositLogs, log)
			}

			if utils.LogIsEvent(log, y.yieldAbi, Withdraw) {
				yieldWithdrawLogs = append(yieldWithdrawLogs, log)
			}

		}
	}

	deposits, err := types.DecodeDepositLog(yieldDepositLogs, y.yieldAbi, Deposit)
	if err != nil {
		return err
	}

	if len(deposits) > 0 {
		err = types.StoreDeposit(tx, deposits, y.raw.Block.Number)
		if err != nil {
			return err
		}
	}

	withdrawals, err := types.DecodeWithdrawLog(yieldWithdrawLogs, y.yieldAbi, Withdraw)
	if err != nil {
		return err
	}

	if len(withdrawals) > 0 {
		err = types.StoreWithdraw(tx, withdrawals, y.raw.Block.Number)
		if err != nil {
			return err
		}
	}

	return nil
}
