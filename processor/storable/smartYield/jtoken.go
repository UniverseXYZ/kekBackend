package smartYield

import (
	"database/sql"
	"encoding/hex"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeJTokenTransferEvent(log web3types.Log, event string) (*types.Transfer, error) {
	var t types.Transfer
	t.TokenAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(&t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.From = utils.Topic2Address(log.Topics[1])
	t.To = utils.Topic2Address(log.Topics[2])
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

func (s *Storable) storeJTokenTransfers(tx *sql.Tx) error {
	if len(s.processed.jTokenTransfers) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("jtoken_transfers", "sy_address", "tx_hash", "tx_index", "log_index", "sender", "receiver", "value", "included_in_block", "block_timestamp"))
	if err != nil {
		return err
	}

	for _, t := range s.processed.jTokenTransfers {
		_, err = stmt.Exec(t.TokenAddress, t.TransactionHash, t.TransactionIndex, t.LogIndex, t.From, t.To, t.Value.String(), s.processed.blockNumber, s.processed.blockTimestamp)
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
