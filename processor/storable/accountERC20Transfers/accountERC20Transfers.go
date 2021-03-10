package accountERC20Transfers

import (
	"database/sql"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type Storable struct {
	config   Config
	raw      *types.RawData
	erc20ABI abi.ABI

	processed struct {
		transfers      []types.Transfer
		blockNumber    int64
		blockTimestamp int64
	}
}

func NewStorable(config Config, raw *types.RawData, erc20ABI abi.ABI) *Storable {
	return &Storable{
		config:   config,
		raw:      raw,
		erc20ABI: erc20ABI,
	}

}

func (s *Storable) ToDB(tx *sql.Tx) error {
	var logs []web3types.Log
	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if utils.LogIsEvent(log, s.erc20ABI, "Transfer") &&
				s.addressExist(log) {
				logs = append(logs, log)
			}
		}
	}

	err := s.decodeLogs(logs)
	if err != nil {
		return err
	}

	s.processed.blockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	s.processed.blockTimestamp, err = strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	err = s.storeTransfers(tx)
	if err != nil {
		return err
	}

	return nil
}
