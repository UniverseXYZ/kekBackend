package supernova

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

var log = logrus.WithField("module", "storable(supernova)")

type SupernovaStorable struct {
	config Config

	Raw          *types.RawData
	supernovaAbi abi.ABI

	Preprocessed struct {
		BlockTimestamp int64
		BlockNumber    int64
	}
}

func NewSupernovaStorable(config Config, raw *types.RawData, supernovaAbi abi.ABI) *SupernovaStorable {
	return &SupernovaStorable{
		config:       config,
		Raw:          raw,
		supernovaAbi: supernovaAbi,
	}
}

func (b *SupernovaStorable) preprocess() error {
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

func (b *SupernovaStorable) ToDB(tx *sql.Tx) error {
	err := b.preprocess()
	if err != nil {
		return err
	}

	var supernovaLogs []web3types.Log
	for _, data := range b.Raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(b.config.SupernovaAddress) {
				supernovaLogs = append(supernovaLogs, log)
			}
		}
	}

	if len(supernovaLogs) == 0 {
		log.Debug("nothing to process")

		return nil
	}

	err = b.handleStakingActions(supernovaLogs, tx)
	if err != nil {
		return err
	}

	err = b.handleLocks(supernovaLogs, tx)
	if err != nil {
		return err
	}

	err = b.handleDelegate(supernovaLogs, tx)
	if err != nil {
		return err
	}

	return nil
}
