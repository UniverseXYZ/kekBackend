package accountERC20Transfers

import (
	"database/sql"
	"encoding/hex"
	"math/big"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeTransfer(log web3types.Log) (*types.Transfer, error) {
	var t types.Transfer
	t.TokenAddress = utils.NormalizeAddress(log.Address)
	t.From = utils.Topic2Address(log.Topics[1])
	t.To = utils.Topic2Address(log.Topics[2])

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}
	var decoded = make(map[string]interface{})

	err = s.erc20ABI.UnpackIntoMap(decoded, "Transfer", data)

	if decoded["value"] != nil {
		t.Value = decoded["value"].(*big.Int)
	} else if decoded["amount"] != nil {
		t.Value = decoded["amount"].(*big.Int)
	} else {
		t.Value = new(big.Int)
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}
	t.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from bond contract to int64")
	}

	t.TransactionHash = log.TransactionHash
	t.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  bond contract to int64")
	}

	return &t, nil
}

func (s *Storable) storeTransfers(tx *sql.Tx) error {
	if len(s.processed.transfers) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("account_erc20_transfers", "token_address", "sender", "receiver", "value", "tx_hash", "tx_index", "log_index", "included_in_block", "block_timestamp"))
	if err != nil {
		return err
	}

	for _, t := range s.processed.transfers {
		_, err = stmt.Exec(t.TokenAddress, t.From, t.To, t.Value.String(), t.TransactionHash, t.TransactionIndex, t.LogIndex, s.processed.blockNumber, s.processed.blockTimestamp)
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
