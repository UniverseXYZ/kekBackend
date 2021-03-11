package smartYield

import (
	"database/sql"
	"encoding/hex"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Harvest struct {
	*types.Event

	ControllerAddress   string
	Caller              string
	CompRewardTotal     *big.Int
	CompRewardSold      *big.Int
	UnderlyingPoolShare *big.Int
	UnderlyingReward    *big.Int
	HarvestCost         *big.Int
}

func (s *Storable) decodeHarvestEvent(log web3types.Log) (*Harvest, error) {
	var h Harvest
	h.ControllerAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["compoundcontroller"].UnpackIntoInterface(&h, HarvestEvent, data)
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

func (s *Storable) storeHarvest(tx *sql.Tx) error {
	if len(s.processed.compoundController.harvests) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("compound_controller_harvests", "controller_address", "caller_address", "comp_reward_total", "comp_reward_sold", "underlying_pool_share", "underlying_reward", "harvest_cost", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.compoundController.harvests {
		_, err = stmt.Exec(a.ControllerAddress, a.Caller, a.CompRewardTotal.String(), a.CompRewardSold.String(), a.UnderlyingPoolShare.String(), a.UnderlyingReward.String(), a.HarvestCost.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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
