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
		tokenActions       TokenTrades
		seniorActions      SeniorTrades
		juniorActions      JuniorTrades
		jTokenTransfers    []types.Transfer
		smartBondTransfers []SmartBondTransfer
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

	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.SmartYieldAddress) {
				smartYieldLogs = append(smartYieldLogs, log)
			}

			if len(log.Topics) == 0 {
				continue
			}

			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.JuniorBondAddress) && utils.LogIsEvent(log, s.abis["juniorbond"], TRANSFER_EVENT) {
				a, err := s.decodeSmartBondTransferEvent(log)
				if err != nil {
					return err
				}
				if a != nil {
					s.processed.smartBondTransfers = append(s.processed.smartBondTransfers, *a)
				}
				continue
			}

			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(s.config.SeniorBondAddress) && utils.LogIsEvent(log, s.abis["seniorbond"], TRANSFER_EVENT) {
				a, err := s.decodeSmartBondTransferEvent(log)
				if err != nil {
					return err
				}
				if a != nil {
					s.processed.smartBondTransfers = append(s.processed.smartBondTransfers, *a)
				}
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

func (s Storable) decodeSmartYieldLog(logs []web3types.Log) error {
	for _, log := range logs {
		if utils.LogIsEvent(log, s.abis["smartyield"], BUY_TOKENS_EVENT) {
			a, err := s.decodeTokenBuyEvent(log, BUY_TOKENS_EVENT)
			if err != nil {
				return err
			}

			if a != nil {
				s.processed.tokenActions.tokenBuyTrades = append(s.processed.tokenActions.tokenBuyTrades, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], SELL_TOKENS_EVENT) {
			a, err := s.decodeTokenSellEvent(log, SELL_TOKENS_EVENT)
			if err != nil {
				return err
			}

			if a != nil {
				s.processed.tokenActions.tokenSellTrades = append(s.processed.tokenActions.tokenSellTrades, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], BUY_SENIOR_BOND_EVENT) {
			a, err := s.decodeSeniorBondBuyEvent(log, BUY_SENIOR_BOND_EVENT)
			if err != nil {
				return err
			}

			if a != nil {
				s.processed.seniorActions.seniorBondBuys = append(s.processed.seniorActions.seniorBondBuys, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], REDEEM_SENIOR_BOND_EVENT) {
			a, err := s.decodeSeniorBondRedeemEvent(log, REDEEM_SENIOR_BOND_EVENT)
			if err != nil {
				return err
			}

			if a != nil {
				s.processed.seniorActions.seniorBondRedeems = append(s.processed.seniorActions.seniorBondRedeems, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], BUY_JUNIOR_BOND_EVENT) {
			a, err := s.decodeJuniorBondBuyEvent(log, BUY_JUNIOR_BOND_EVENT)
			if err != nil {
				return err
			}

			if a != nil {
				s.processed.juniorActions.juniorBondBuys = append(s.processed.juniorActions.juniorBondBuys, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], REDEEM_JUNIOR_BOND_EVENT) {
			a, err := s.decodeJuniorBondRedeemEvent(log, REDEEM_JUNIOR_BOND_EVENT)
			if err != nil {
				return err
			}

			if a != nil {
				s.processed.juniorActions.juniorBondRedeems = append(s.processed.juniorActions.juniorBondRedeems, *a)
			}
			continue
		}
		if utils.LogIsEvent(log, s.abis["smartyield"], TRANSFER_EVENT) {
			a, err := s.decodeSmartYieldTransferEvent(log, TRANSFER_EVENT)
			if err != nil {
				return err
			}
			if a != nil {
				s.processed.jTokenTransfers = append(s.processed.jTokenTransfers, *a)
			}
			continue
		}
	}

	return nil
}
