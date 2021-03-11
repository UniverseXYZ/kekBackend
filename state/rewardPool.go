package state

import (
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func RewardPools() []types.SYRewardPool {
	return instance.rewardPools
}

func RewardPoolByAddress(address string) *types.SYRewardPool {
	for _, p := range instance.rewardPools {
		if utils.NormalizeAddress(address) == utils.NormalizeAddress(p.PoolAddress) {
			return &p
		}
	}

	return nil
}
