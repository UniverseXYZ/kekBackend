package smartYield

import (
	"database/sql"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

var log = logrus.WithField("module", "storable(smart yield)")

type Storable struct {
	config Config
	raw    *types.RawData
	abis   map[string]abi.ABI

	processed struct {
		tokenActions     TokenTrades
		seniorActions    SeniorTrades
		juniorActions    JuniorTrades
		jTokenTransfers  []types.Transfer
		sTokenTransfers  []STokenTransfer
		compoundProvider CompoundProvider
		blockNumber      int64
		blockTimestamp   int64
	}
}

func NewStorable(config Config, raw *types.RawData, abis map[string]abi.ABI) *Storable {
	return &Storable{
		config: config,
		raw:    raw,
		abis:   abis,
	}
}

func (s *Storable) ToDB(tx *sql.Tx) error {
	var smartYieldLogs []web3types.Log
	var compoundProviderLogs []web3types.Log

	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.SmartYieldAddress) {
				smartYieldLogs = append(smartYieldLogs, log)
			}

			if len(log.Topics) == 0 {
				continue
			}

			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.JuniorBondAddress) && utils.LogIsEvent(log, s.abis["juniorbond"], TRANSFER_EVENT) {
				a, err := s.decodeSTokenTransferEvent(log)
				if err != nil {
					return err
				}
				if a != nil {
					s.processed.sTokenTransfers = append(s.processed.sTokenTransfers, *a)
				}
				continue
			}

			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.SeniorBondAddress) && utils.LogIsEvent(log, s.abis["seniorbond"], TRANSFER_EVENT) {
				a, err := s.decodeSTokenTransferEvent(log)
				if err != nil {
					return err
				}
				if a != nil {
					s.processed.sTokenTransfers = append(s.processed.sTokenTransfers, *a)
				}
				continue
			}

			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.CompoundProviderAddress) {
				compoundProviderLogs = append(compoundProviderLogs, log)
				continue
			}

		}
	}

	if len(smartYieldLogs) == 0 {
		log.WithField("handler", "smart yield trades").Debug("no actions found")
		return nil
	}

	err := s.decodeSmartYieldLog(smartYieldLogs)
	if err != nil {
		return err
	}

	err = s.decodeCompoundProviderEvents(compoundProviderLogs)
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

	err = s.storeProcessed(tx)
	if err != nil {
		return err
	}

	return nil
}
