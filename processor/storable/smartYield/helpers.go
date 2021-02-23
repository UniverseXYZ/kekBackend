package smartYield

import (
	web3types "github.com/alethio/web3-go/types"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeSmartYieldLog(logs []web3types.Log) error {
	for _, log := range logs {
		if utils.LogIsEvent(log, s.abis["smartyield"], BuyTokensEvent) {
			a, err := s.decodeTokenBuyEvent(log, BuyTokensEvent)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.tokenActions.tokenBuyTrades = append(s.processed.tokenActions.tokenBuyTrades, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], SellTokensEvent) {
			a, err := s.decodeTokenSellEvent(log, SellTokensEvent)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.tokenActions.tokenSellTrades = append(s.processed.tokenActions.tokenSellTrades, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], BuySeniorBondEvent) {
			a, err := s.decodeSeniorBondBuyEvent(log, BuySeniorBondEvent)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.seniorActions.seniorBondBuys = append(s.processed.seniorActions.seniorBondBuys, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], RedeemSeniorBondEvent) {
			a, err := s.decodeSeniorBondRedeemEvent(log, RedeemSeniorBondEvent)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.seniorActions.seniorBondRedeems = append(s.processed.seniorActions.seniorBondRedeems, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], BuyJuniorBondEvent) {
			a, err := s.decodeJuniorBondBuyEvent(log, BuyJuniorBondEvent)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.juniorActions.juniorBondBuys = append(s.processed.juniorActions.juniorBondBuys, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], RedeemJuniorBondEvent) {
			a, err := s.decodeJuniorBondRedeemEvent(log, RedeemJuniorBondEvent)
			if err != nil {
				return err
			} else if a != nil {
				s.processed.juniorActions.juniorBondRedeems = append(s.processed.juniorActions.juniorBondRedeems, *a)
			}
			continue
		}

		if utils.LogIsEvent(log, s.abis["smartyield"], TransferEvent) {
			a, err := s.decodeJTokenTransferEvent(log, TransferEvent)
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
		if utils.LogIsEvent(log, s.abis["compoundprovider"], HarvestEvent) {
			a, err := s.decodeHarvestEvent(log)
			if err != nil {
				return err
			}
			if a != nil {
				s.processed.compoundProvider.harvests = append(s.processed.compoundProvider.harvests, *a)
			}
		}

		if utils.LogIsEvent(log, s.abis["compoundprovider"], TransferFeesEvent) {
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
