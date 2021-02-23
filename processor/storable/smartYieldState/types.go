package smartYieldState

import (
	"github.com/shopspring/decimal"
)

type State struct {
	PoolAddress string

	TotalLiquidity  decimal.Decimal
	JuniorLiquidity decimal.Decimal
	JTokenPrice     decimal.Decimal

	SeniorAPY        float64
	JuniorAPY        float64
	OriginatorApy    float64
	OriginatorNetApy float64
}
