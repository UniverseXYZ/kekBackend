package smartYieldRewards

import (
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeEvents(logs []web3types.Log) error {
	for _, log := range logs {
		if utils.LogIsEvent(log, s.syRewardABI, "Claim") {
			a, err := s.decodeClaimEvent(log, "Claim")
			if err != nil {
				return err
			}

			if a != nil {
				p := state.PoolByAddress(log.Address)
				a.PoolAddress = p.PoolAddress
				s.processed.claims = append(s.processed.claims, *a)
				continue
			}
		}

		if utils.LogIsEvent(log, s.syRewardABI, "Deposit") {
			a, err := s.decodeStakingEvent(log)
			if err != nil {
				return errors.Wrap(err, "coud not decode deposit event")
			}

			if a != nil {
				p := state.PoolByAddress(log.Address)
				a.PoolAddress = p.PoolAddress
				a.ActionType = DEPOSIT
				s.processed.stakingActions = append(s.processed.stakingActions, *a)
				continue
			}

		}

		if utils.LogIsEvent(log, s.syRewardABI, "Withdraw") {
			a, err := s.decodeStakingEvent(log)
			if err != nil {
				return errors.Wrap(err, "coud not decode withdraw event")
			}

			if a != nil {
				p := state.PoolByAddress(log.Address)
				a.PoolAddress = p.PoolAddress
				a.ActionType = WITHDRAW
				s.processed.stakingActions = append(s.processed.stakingActions, *a)
				continue
			}
		}
	}

	return nil
}

func (s *Storable) decodeClaimEvent(log web3types.Log, action string) (*ClaimEvent, error) {
	var e ClaimEvent
	var err error

	e.UserAddress = utils.Topic2Address(log.Topics[1])
	e.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.syRewardABI.UnpackIntoInterface(&e, action, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	return &e, nil
}

func (s *Storable) decodeStakingEvent(log web3types.Log) (*StakingAction, error) {
	return nil, nil
}
