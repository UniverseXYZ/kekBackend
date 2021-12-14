package accountERC20Transfers

import (
	"database/sql"
	"strconv"

	"github.com/alethio/web3-go/ethrpc"
	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/state"
	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

type Storable struct {
	config   Config
	raw      *types.RawData
	erc20ABI abi.ABI
	ethConn  *ethclient.Client

	processed struct {
		transfers      []types.Transfer
		blockNumber    int64
		blockTimestamp int64
	}
}

func NewStorable(config Config, raw *types.RawData, erc20ABI abi.ABI, ethConn *ethclient.Client) *Storable {
	return &Storable{
		config:   config,
		raw:      raw,
		erc20ABI: erc20ABI,
		ethConn:  ethConn,
	}

}

func (s *Storable) ToDB(tx *sql.Tx, ethBatch *ethrpc.ETH) error {
	var logs []web3types.Log
	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if utils.LogIsEvent(log, s.erc20ABI, "Transfer") && state.IsMonitoredAccount(log) {
				err := s.checkTokenExists(tx, utils.NormalizeAddress(log.Address))
				if err != nil {
					continue
				}

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
