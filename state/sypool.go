package state

import (
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func Pools() []types.SYPool {
	return instance.pools
}

func PoolBySmartYieldAddress(address string) *types.SYPool {
	for _, p := range instance.pools {
		if utils.NormalizeAddress(address) == utils.NormalizeAddress(p.SmartYieldAddress) {
			return &p
		}
	}

	return nil
}

func PoolByJuniorBondAddress(address string) *types.SYPool {
	for _, p := range instance.pools {
		if utils.NormalizeAddress(address) == utils.NormalizeAddress(p.JuniorBondAddress) {
			return &p
		}
	}

	return nil
}

func PoolBySeniorBondAddress(address string) *types.SYPool {
	for _, p := range instance.pools {
		if utils.NormalizeAddress(address) == utils.NormalizeAddress(p.SeniorBondAddress) {
			return &p
		}
	}

	return nil
}

func PoolByProviderAddress(address string) *types.SYPool {
	for _, p := range instance.pools {
		if utils.NormalizeAddress(address) == utils.NormalizeAddress(p.ProviderAddress) {
			return &p
		}
	}

	return nil
}

func PoolByUnderlyingAddress(address string) *types.SYPool {
	for _, p := range instance.pools {
		if utils.NormalizeAddress(address) == utils.NormalizeAddress(p.UnderlyingAddress) {
			return &p
		}
	}

	return nil
}
