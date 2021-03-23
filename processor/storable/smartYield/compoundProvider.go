package smartYield

import (
	"database/sql"
	"encoding/hex"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type TransferFees struct {
	*types.Event

	ProviderAddress string
	Caller          string
	FeesOwner       string
	Fees            *big.Int
}

func (s *Storable) decodeTransferFeesEvent(log web3types.Log) (*TransferFees, error) {
	var t TransferFees
	t.ProviderAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["compoundprovider"].UnpackIntoInterface(&t, TransferFeesEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.Caller = utils.Topic2Address(log.Topics[1])
	t.FeesOwner = utils.Topic2Address(log.Topics[2])
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) storeTransferFees(tx *sql.Tx) error {
	if len(s.processed.compoundProvider.transfersFees) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("compound_provider_transfer_fees", "provider_address", "caller_address", "fees_owner", "fees", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.compoundProvider.transfersFees {
		_, err = stmt.Exec(a.ProviderAddress, a.Caller, a.FeesOwner, a.Fees.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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
