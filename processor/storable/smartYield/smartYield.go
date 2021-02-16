package smartYield

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type Storable struct {
	config Config
	raw    *types.RawData
	abis   map[string]abi.ABI

	processed struct {
		tokenActions  TokenTrades
		seniorActions SeniorTrades
		juniorActions JuniorTrades
	}
}

func NewStorable(config Config, raw *types.RawData, abis map[string]abi.ABI) (*Storable, error) {
	if _, exist := abis["smartyield"]; !exist {
		return nil, errors.New("could not find smartYield abi")
	}

	if _, exist := abis["erc20"]; !exist {
		return nil, errors.New("could not find erc20 token abi")
	}

	if _, exist := abis["erc721"]; !exist {
		return nil, errors.New("could not find erc20 token abi")
	}

	return &Storable{
		config: config,
		raw:    raw,
		abis:   abis,
	}, nil
}

func (s *Storable) ToDB(tx *sql.Tx) error {
	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(s.config.Address) {
				continue
			}

			if len(log.Topics) == 0 {
				continue
			}

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
		}

	}

	return nil
}
