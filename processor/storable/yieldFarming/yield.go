package yieldFarming

import (
	"database/sql"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/alethio/web3-go/ethrpc"
	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

var log = logrus.WithField("module", "storable(yield farming)")

type Storable struct {
	config   Config
	raw      *types.RawData
	yieldAbi abi.ABI
}

type StakingAction struct {
	UserAddress      string
	TokenAddress     string
	Amount           *big.Int
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
	ActionType       ActionType
}

func NewStorable(config Config, raw *types.RawData, yieldAbi abi.ABI) *Storable {
	return &Storable{
		config:   config,
		raw:      raw,
		yieldAbi: yieldAbi,
	}
}

func (y *Storable) ToDB(tx *sql.Tx, ethBatch *ethrpc.ETH) error {
	var stakingActions []StakingAction

	for _, data := range y.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(y.config.Address) {
				continue
			}
			if len(log.Topics) == 0 {
				continue
			}

			if utils.LogIsEvent(log, y.yieldAbi, Deposit) {
				d, err := y.decodeLog(log, Deposit)
				if err != nil {
					return err
				}

				d.ActionType = DEPOSIT
				stakingActions = append(stakingActions, *d)
			}

			if utils.LogIsEvent(log, y.yieldAbi, Withdraw) {
				w, err := y.decodeLog(log, Withdraw)
				if err != nil {
					return err
				}

				w.ActionType = WITHDRAW
				stakingActions = append(stakingActions, *w)
			}

		}
	}
	if len(stakingActions) == 0 {
		log.WithField("handler", "staking actions").Debug("no actions found")
		return nil
	}

	err := y.storeActions(tx, stakingActions)
	if err != nil {
		return err
	}

	return nil
}

func (y Storable) decodeLog(log web3types.Log, event string) (*StakingAction, error) {
	var d StakingAction

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	var extraData = make(map[string]interface{})
	err = y.yieldAbi.UnpackIntoMap(extraData, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	d.UserAddress = utils.Topic2Address(log.Topics[1])
	d.TokenAddress = utils.Topic2Address(log.Topics[2])
	d.Amount = extraData["amount"].(*big.Int)
	d.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from kek contract to int64")
	}

	d.TransactionHash = log.TransactionHash
	d.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  kek contract to int64")
	}

	return &d, nil
}

func (y Storable) storeActions(tx *sql.Tx, actions []StakingAction) error {
	stmt, err := tx.Prepare(pq.CopyIn("yield_farming_actions", "tx_hash", "tx_index", "log_index", "user_address", "token_address", "amount", "action_type", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(y.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(y.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.UserAddress, a.TokenAddress, a.Amount.String(), a.ActionType, blockTimestamp, blockNumber)
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
