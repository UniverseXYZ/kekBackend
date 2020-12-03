package bond

import (
	"database/sql"
	"encoding/hex"
	"math/big"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type BondStorable struct {
	config  Config
	Raw     *types.RawData
	bondAbi abi.ABI
}
type Transfer struct {
	From             string
	To               string
	Value            *big.Int
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64

	//from ,to(String) ,value(big.int) maybe pointers
	//txHash(string)
	//txIndex logIndex(int)
}

func NewBondStorable(config Config, raw *types.RawData, bondAbi abi.ABI) *BondStorable {
	return &BondStorable{
		config:  config,
		Raw:     raw,
		bondAbi: bondAbi,
	}
}

func (b BondStorable) ToDB(tx *sql.Tx) error {
	var bondTransfers []web3types.Log
	var transfers []Transfer
	for _, data := range b.Raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(b.config.BondAddress) {
				continue
			}

			if len(log.Topics) == 0 {
				continue
			}

			if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(b.bondAbi.Events["Transfer"].ID.String()) {
				bondTransfers = append(bondTransfers, log)
			}
		}
	}

	for _, log := range bondTransfers {
		var t Transfer
		data, err := hex.DecodeString(utils.Trim0x(log.Data))
		if err != nil {
			return errors.Wrap(err, "could not decode log data")
		}

		err = b.bondAbi.UnpackIntoInterface(&t, "Transfer", data)
		if err != nil {
			return errors.Wrap(err, "could not unpack log data")
		}

		t.From = utils.Topic2Address(log.Topics[1])
		t.To = utils.Topic2Address(log.Topics[2])
		t.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
		if err != nil {
			return errors.Wrap(err, "could not convert transactionIndex to int64")
		}

		t.TransactionHash = log.TransactionHash
		t.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
		if err != nil {
			return errors.Wrap(err, "could not convert logIndex to int64")
		}

		transfers = append(transfers, t)
	}
	spew.Dump(transfers)

	stmt, err := tx.Prepare(pq.CopyIn("bond_transfers", "tx_hash", "tx_index", "log_index", "sender", "receiver", "value", "included_in_block"))
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
