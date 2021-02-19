package smartYield

import (
	web3types "github.com/alethio/web3-go/types"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeSmartYieldLog(logs []web3types.Log) error {
	for _, log := range logs {
		if utils.LogIsEvent(log, s.abis["smartyield"], BUY_TOKENS_EVENT) {
			a, err := s.decodeTokenBuyEvent(log, BUY_TOKENS_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.tokenActions.tokenBuyTrades = append(s.processed.tokenActions.tokenBuyTrades, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], SELL_TOKENS_EVENT) {
			a, err := s.decodeTokenSellEvent(log, SELL_TOKENS_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.tokenActions.tokenSellTrades = append(s.processed.tokenActions.tokenSellTrades, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], BUY_SENIOR_BOND_EVENT) {
			a, err := s.decodeSeniorBondBuyEvent(log, BUY_SENIOR_BOND_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.seniorActions.seniorBondBuys = append(s.processed.seniorActions.seniorBondBuys, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], REDEEM_SENIOR_BOND_EVENT) {
			a, err := s.decodeSeniorBondRedeemEvent(log, REDEEM_SENIOR_BOND_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.seniorActions.seniorBondRedeems = append(s.processed.seniorActions.seniorBondRedeems, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], BUY_JUNIOR_BOND_EVENT) {
			a, err := s.decodeJuniorBondBuyEvent(log, BUY_JUNIOR_BOND_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.juniorActions.juniorBondBuys = append(s.processed.juniorActions.juniorBondBuys, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], REDEEM_JUNIOR_BOND_EVENT) {
			a, err := s.decodeJuniorBondRedeemEvent(log, REDEEM_JUNIOR_BOND_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.juniorActions.juniorBondRedeems = append(s.processed.juniorActions.juniorBondRedeems, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], TRANSFER_EVENT) {
			a, err := s.decodeJTokenTransferEvent(log, TRANSFER_EVENT)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.jTokenTransfers = append(s.processed.jTokenTransfers, *a)
			}
			continue
		}
	}

	return nil
}

func (s *Storable) decodeCompoundProviderEvents(logs []web3types.Log) error {
	for _, log := range logs {
		if utils.LogIsEvent(log, s.abis["compoundprovider"], HARVEST_EVENT) {
			a, err := s.decodeHarvestEvent(log)
			if err != nil {
				return err
			}
			if a != nil {
				s.processed.compoundProvider.harvests = append(s.processed.compoundProvider.harvests, *a)
			}
		}

		if utils.LogIsEvent(log, s.abis["compoundprovider"], TRANSFER_FEES_EVENT) {
			a, err := s.decodeTransferFeesEvent(log)
			if err != nil {
				return err
			}
			if a != nil {
				s.processed.compoundProvider.transfersFees = append(s.processed.compoundProvider.transfersFees, *a)
			}
		}
	}

	return nil
}
