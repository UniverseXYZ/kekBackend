package smartYieldState

import "math/big"

type State struct {
	PoolAddress string

	TotalLiquidity  *big.Int
	JuniorLiquidity *big.Int
	JTokenPrice     *big.Int

	SeniorAPY        float64
	JuniorAPY        float64
	OriginatorApy    float64
	OriginatorNetApy float64
}
