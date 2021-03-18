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

func AddNewPoolToState(pool types.SYRewardPool) {
	instance.rewardPools = append(instance.rewardPools, pool)
}

func AddNewPoolToDB(pool types.SYRewardPool) error {
	_, err := instance.db.Exec(`insert into smart_yield_reward_pools (pool_address, pool_token_address, reward_token_address,start_at_block) values ($1,$2,$3,$4)`,
		pool.PoolAddress, pool.PoolTokenAddress, pool.RewardTokenAddress, pool.StartAtBlock)

	if err != nil {
		return err
	}

	return nil
}
