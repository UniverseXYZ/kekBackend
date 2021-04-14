package barn

import (
	"database/sql"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

var log = logrus.WithField("module", "storable(barn)")

type BarnStorable struct {
	config Config

	Raw     *types.RawData
	barnAbi abi.ABI

	Preprocessed struct {
		BlockTimestamp int64
		BlockNumber    int64
	}
}

func NewBarnStorable(config Config, raw *types.RawData, barnAbi abi.ABI) *BarnStorable {
	return &BarnStorable{
		config:  config,
		Raw:     raw,
		barnAbi: barnAbi,
	}
}

func (b *BarnStorable) preprocess() error {
	var err error

	b.Preprocessed.BlockNumber, err = strconv.ParseInt(b.Raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "unable to process block number")
	}

	b.Preprocessed.BlockTimestamp, err = strconv.ParseInt(b.Raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse block timestamp")
	}

	return nil
}

func (b *BarnStorable) ToDB(tx *sql.Tx) error {
	err := b.preprocess()
	if err != nil {
		return err
	}

	var barnLogs []web3types.Log
	for _, data := range b.Raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(b.config.BarnAddress) {
				barnLogs = append(barnLogs, log)
			}
		}
	}

	if len(barnLogs) == 0 {
		log.Debug("nothing to process")

		return nil
	}

	err = b.handleStakingActions(barnLogs, tx)
	if err != nil {
		return err
	}

	err = b.handleLocks(barnLogs, tx)
	if err != nil {
		return err
	}

	err = b.handleDelegate(barnLogs, tx)
	if err != nil {
		return err
	}

	return nil
}
