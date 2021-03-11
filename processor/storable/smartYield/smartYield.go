package smartYield

import (
	"database/sql"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

var log = logrus.WithField("module", "storable(smart yield)")

type Storable struct {
	config Config
	raw    *types.RawData
	abis   map[string]abi.ABI

	processed struct {
		tokenActions       TokenTrades
		seniorActions      SeniorTrades
		juniorActions      JuniorTrades
		jTokenTransfers    []types.Transfer
		ERC721Transfers    []ERC721Transfer
		compoundProvider   CompoundProvider
		compoundController CompoundController
		blockNumber        int64
		blockTimestamp     int64
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
	var compoundControllerLogs []web3types.Log

	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if state.PoolBySmartYieldAddress(log.Address) != nil {
				smartYieldLogs = append(smartYieldLogs, log)
				continue
			}

			if state.PoolByJuniorBondAddress(log.Address) != nil && utils.LogIsEvent(log, s.abis["juniorbond"], TransferEvent) {
				p := state.PoolByJuniorBondAddress(log.Address)
				a, err := s.decodeERC721TransferEvent(log)

				if err != nil {
					return err
				} else if a != nil {
					a.ProtocolId = p.ProtocolId
					a.SYAddress = p.SmartYieldAddress
					a.TokenType = "junior"
					s.processed.ERC721Transfers = append(s.processed.ERC721Transfers, *a)
				}
				continue
			}

			if state.PoolBySeniorBondAddress(log.Address) != nil && utils.LogIsEvent(log, s.abis["seniorbond"], TransferEvent) {
				p := state.PoolBySeniorBondAddress(log.Address)
				a, err := s.decodeERC721TransferEvent(log)
				if err != nil {
					return err
				} else if a != nil {
					a.ProtocolId = p.ProtocolId
					a.SYAddress = p.SmartYieldAddress
					a.TokenType = "senior"
					s.processed.ERC721Transfers = append(s.processed.ERC721Transfers, *a)
				}
				continue
			}

			if state.PoolByProviderAddress(log.Address) != nil {
				compoundProviderLogs = append(compoundProviderLogs, log)
				continue
			}

			if state.PoolByControllerAddress(log.Address) != nil {
				compoundControllerLogs = append(compoundControllerLogs, log)
				continue
			}
		}
	}

	err := s.decodeSmartYieldLog(smartYieldLogs)
	if err != nil {
		return err
	}

	err = s.decodeCompoundProviderEvents(compoundProviderLogs)
	if err != nil {
		return err
	}

	err = s.decodeCompoundControllerEvents(compoundControllerLogs)
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
