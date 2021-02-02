package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProposalState string

const (
	CREATED   ProposalState = "CREATED"
	WARMUP    ProposalState = "WARMUP"
	ACTIVE    ProposalState = "ACTIVE"
	CANCELED  ProposalState = "CANCELED"
	FAILED    ProposalState = "FAILED"
	ACCEPTED  ProposalState = "ACCEPTED"
	QUEUED    ProposalState = "QUEUED"
	GRACE     ProposalState = "GRACE"
	EXPIRED   ProposalState = "EXPIRED"
	EXECUTED  ProposalState = "EXECUTED"
	ABROGATED ProposalState = "ABROGATED"
)
const epochDuration = 604800
const maxEpoch = 100

var epoch1StartUnix = int64(1603065600)
var epoch1Start = time.Unix(epoch1StartUnix, 0).UTC()
var vestingTotal = decimal.New(2200000, 0)
var totalSupply = decimal.New(10000000, 0)

var decimals = map[string]int32{
	"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": 6,
	"0x57ab1ec28d129707052df4df418d58a2d46d5f51": 18,
	"0x6b175474e89094c44da98b954eedeac495271d0f": 18,
	"0x6591c4bcd6d7a1eb4e537da8b78676c1576ba244": 18,
	"0x0391d2021f89dc339f60fff84546ea23e337750f": 18,
}

var pools = map[string]Pool{
	"stable": {
		tokens: []string{
			"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			"0x57Ab1ec28D129707052df4dF418D58a2D46d5f51",
			"0x6B175474E89094C44Da98b954EedeAC495271d0F",
		},
		epochDelayFromStaking: 0,
	},
	"unilp": {
		tokens:                []string{"0x6591c4BcD6D7A1eb4E537DA8B78676C1576Ba244"},
		epochDelayFromStaking: 1,
	},
	"bond": {
		tokens:                []string{"0x0391D2021f89DC339F60Fff84546EA23E337750f"},
		epochDelayFromStaking: 4,
	},
}
