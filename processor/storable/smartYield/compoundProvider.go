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

type Harvest struct {
	*types.Event
	Caller              string
	UnderlyingGot       *big.Int
	RewardExpected      *big.Int
	UnderlyingDeposited *big.Int
	Fees                *big.Int
	Reward              *big.Int
}

type TransferFees struct {
	*types.Event
	Caller    string
	FeesOwner string
	Fees      *big.Int
}

func (s *Storable) decodeHarvestEvent(log web3types.Log) (*Harvest, error) {
	var h Harvest
	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["compoundprovider"].UnpackIntoInterface(&h, HARVEST_EVENT, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	h.Caller = utils.Topic2Address(log.Topics[1])
	h.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *Storable) decodeTransferFeesEvent(log web3types.Log) (*TransferFees, error) {
	var t TransferFees

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["compoundprovider"].UnpackIntoInterface(&t, TRANSFER_FEES_EVENT, data)
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

func (s *Storable) storeHarvest(tx *sql.Tx) error {
	if len(s.processed.compoundProvider.harvests) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("compound_provider_harvest", "caller_address", "underlying_got", "reward_expected", "underlying_deposited", "fees", "reward", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.compoundProvider.harvests {
		_, err = stmt.Exec(a.Caller, a.UnderlyingGot.String(), a.RewardExpected.String(), a.UnderlyingDeposited.String(), a.Fees.String(), a.Reward.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

func (s *Storable) storeTransferFees(tx *sql.Tx) error {
	if len(s.processed.compoundProvider.transfersFees) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("compound_provider_transfer_fees", "caller_address", "fees_owner", "fees", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.compoundProvider.transfersFees {
		_, err = stmt.Exec(a.Caller, a.FeesOwner, a.Fees.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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
