package kek

import (
	"database/sql"
	"encoding/hex"
	"strconv"

	"github.com/alethio/web3-go/ethrpc"
	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

type KekStorable struct {
	config Config
	Raw    *types.RawData
	kekAbi abi.ABI
}

func NewKekStorable(config Config, raw *types.RawData, kekAbi abi.ABI) *KekStorable {
	return &KekStorable{
		config: config,
		Raw:    raw,
		kekAbi: kekAbi,
	}
}

func (b KekStorable) ToDB(tx *sql.Tx, ethBatch *ethrpc.ETH) error {
	var kekTransfers []web3types.Log
	var transfers []types.Transfer
	for _, data := range b.Raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(b.config.KekAddress) {
				continue
			}

			if len(log.Topics) == 0 {
				continue
			}

			if utils.LogIsEvent(log, b.kekAbi, "Transfer") {
				kekTransfers = append(kekTransfers, log)
			}
		}
	}

	for _, log := range kekTransfers {
		var t types.Transfer
		data, err := hex.DecodeString(utils.Trim0x(log.Data))
		if err != nil {
			return errors.Wrap(err, "could not decode log data")
		}

		err = b.kekAbi.UnpackIntoInterface(&t, "Transfer", data)
		if err != nil {
			return errors.Wrap(err, "could not unpack log data")
		}

		t.From = utils.Topic2Address(log.Topics[1])
		t.To = utils.Topic2Address(log.Topics[2])
		t.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
		if err != nil {
			return errors.Wrap(err, "could not convert transactionIndex from kek contract to int64")
		}

		t.TransactionHash = log.TransactionHash
		t.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
		if err != nil {
			return errors.Wrap(err, "could not convert logIndex from  kek contract to int64")
		}

		transfers = append(transfers, t)
	}

	stmt, err := tx.Prepare(pq.CopyIn("kek_transfers", "tx_hash", "tx_index", "log_index", "sender", "receiver", "value", "included_in_block"))
	if err != nil {
		return err
	}

	number, err := strconv.ParseInt(b.Raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, t := range transfers {
		_, err = stmt.Exec(t.TransactionHash, t.TransactionIndex, t.LogIndex, t.From, t.To, t.Value.String(), number)
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
