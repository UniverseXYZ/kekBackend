package smartYieldRewards

import (
	web3types "github.com/alethio/web3-go/types"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeEvents(logs []web3types.Log) error {
	for _, log := range logs {
		if utils.LogIsEvent(log, s.syRewardABI, "Claim") {
			a, err := s.decodeClaimEvent(log)
			if err != nil {
				return err
			}
			if a != nil {
				p := state.PoolByAddress(log.Address)
				a.PoolAddress = p.PoolAddress
				s.processed.claims = append(s.processed.claims, *a)
			}
		}

		if utils.LogIsEvent(log, s.syRewardABI, "Deposit") {

		}
	}

	return nil
}

func (s *Storable) decodeClaimEvent(log web3types.Log) (*ClaimEvent, error) {
	return nil, nil
}

func (s *Storable) decodeStakingEvent(log web3types.Log) (*StakingAction, error) {
	return nil, nil
}
