package types

import (
	"database/sql"
	"encoding/hex"
	"math/big"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

type Withdraw struct {
	UserAddress      string
	TokenAddress     string
	Amount           *big.Int
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

func DecodeWithdrawLog(logs []web3types.Log, abi abi.ABI, event string) ([]Withdraw, error) {
	if len(logs) == 0 {
		return nil, nil
	}

	var withdraws []Withdraw
	for _, log := range logs {
		var w Withdraw

		data, err := hex.DecodeString(utils.Trim0x(log.Data))
		if err != nil {
			return nil, errors.Wrap(err, "could not decode log data")
		}

		var withdrawData = make(map[string]interface{})
		err = abi.UnpackIntoMap(withdrawData, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}

		w.UserAddress = utils.Topic2Address(log.Topics[1])
		w.TokenAddress = utils.Topic2Address(log.Topics[2])
		w.Amount = withdrawData["amount"].(*big.Int)
		w.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert transactionIndex from bond contract to int64")
		}

		w.TransactionHash = log.TransactionHash
		w.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert logIndex from  bond contract to int64")
		}

		withdraws = append(withdraws, w)
	}

	return withdraws, nil
}

func StoreWithdraw(tx *sql.Tx, withdrawals []Withdraw, blockNumber string) error {
	stmt, err := tx.Prepare(pq.CopyIn("withdrawals", "tx_hash", "tx_index", "log_index", "user_address", "token_address", "amount", "included_in_block"))
	if err != nil {
		return err
	}

	number, err := strconv.ParseInt(blockNumber, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, w := range withdrawals {
		_, err = stmt.Exec(w.TransactionHash, w.TransactionIndex, w.LogIndex, w.UserAddress, w.TokenAddress, w.Amount.String(), number)
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
