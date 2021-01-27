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

type Deposit struct {
	UserAddress      string
	TokenAddress     string
	Amount           *big.Int
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

func DecodeDepositLog(logs []web3types.Log, abi abi.ABI, event string) ([]Deposit, error) {
	if len(logs) == 0 {
		return nil, nil
	}

	var deposits []Deposit
	for _, log := range logs {
		var d Deposit

		data, err := hex.DecodeString(utils.Trim0x(log.Data))
		if err != nil {
			return nil, errors.Wrap(err, "could not decode log data")
		}

		var depositData = make(map[string]interface{})
		err = abi.UnpackIntoMap(depositData, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}

		d.UserAddress = utils.Topic2Address(log.Topics[1])
		d.TokenAddress = utils.Topic2Address(log.Topics[2])
		d.Amount = depositData["amount"].(*big.Int)
		d.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert transactionIndex from bond contract to int64")
		}

		d.TransactionHash = log.TransactionHash
		d.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert logIndex from  bond contract to int64")
		}

		deposits = append(deposits, d)
	}

	return deposits, nil
}

func StoreDeposit(tx *sql.Tx, deposits []Deposit, blockNumber string) error {
	stmt, err := tx.Prepare(pq.CopyIn("deposits", "tx_hash", "tx_index", "log_index", "user_address", "token_address", "amount", "included_in_block"))
	if err != nil {
		return err
	}

	number, err := strconv.ParseInt(blockNumber, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, d := range deposits {
		_, err = stmt.Exec(d.TransactionHash, d.TransactionIndex, d.LogIndex, d.UserAddress, d.TokenAddress, d.Amount.String(), number)
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
