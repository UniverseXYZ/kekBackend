package smartYieldState

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
)

type Config struct {
	ComptrollerAddress string
	BlocksPerMinute    int64
	StartAt            int64
}

type Storable struct {
	config Config
	raw    *types.RawData

	abis map[string]abi.ABI
	eth  *ethrpc.ETH

	Preprocessed struct {
		BlockTimestamp time.Time
		BlockNumber    int64
	}
}

func New(config Config, raw *types.RawData, abis map[string]abi.ABI, eth *ethrpc.ETH) (*Storable, error) {
	var s = &Storable{
		config: config,
		raw:    raw,
		abis:   abis,
		eth:    eth,
	}

	var err error
	s.Preprocessed.BlockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process block number")
	}

	txUnix, err := strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse block timestamp")
	}

	s.Preprocessed.BlockTimestamp = time.Unix(txUnix, 0)

	return s, nil
}

func (s Storable) ToDB(tx *sql.Tx) error {
	if s.Preprocessed.BlockNumber < s.config.StartAt {
		return nil
	}

	var wg = &errgroup.Group{}

	var results = make(map[string]*State)
	var mu = &sync.Mutex{}

	for _, p := range state.Pools() {
		p := p

		results[p.SmartYieldAddress] = &State{
			PoolAddress: p.SmartYieldAddress,
		}

		s.getTotalLiquidity(wg, p, mu, results)
		s.getJuniorLiquidity(wg, p, mu, results)
		s.getPrice(wg, p, mu, results)
		s.getMaxBondDailyRate(wg, p, mu, results)
		s.getAbond(wg, p, mu, results)

		if p.ProtocolId == "compound/v2" {
			s.getCompoundAPY(wg, p, mu, results)
		}
	}

	err := wg.Wait()
	if err != nil {
		return err
	}

	for _, p := range state.Pools() {
		// 	(seniorLiq - (abondAPY / originatorAPY * seniorLiq) + juniorLiq) * originatorAPY / juniorLiq
		r := results[p.SmartYieldAddress]

		if r.OriginatorNetApy == 0 || r.JuniorLiquidity.Equal(decimal.NewFromInt(0)) {
			results[p.SmartYieldAddress].JuniorAPY = 0
			continue
		}

		seniorLiq := r.TotalLiquidity.Sub(r.JuniorLiquidity)

		abondGain := decimal.NewFromBigInt(r.Abond.Gain, -int32(p.UnderlyingDecimals))
		abondPrincipal := decimal.NewFromBigInt(r.Abond.Principal, -int32(p.UnderlyingDecimals))
		abondIssuedAt := decimal.NewFromBigInt(r.Abond.IssuedAt, -18)
		abondMaturesAt := decimal.NewFromBigInt(r.Abond.MaturesAt, -18)

		var abondAPY float64
		if !abondPrincipal.Equal(decimal.NewFromInt(0)) {
			abondAPY, _ = abondGain.Div(abondPrincipal).Div(abondMaturesAt.Sub(abondIssuedAt)).Mul(decimal.NewFromInt(365 * 24 * 60 * 60)).Float64()
		}

		a := decimal.NewFromFloat(abondAPY / r.OriginatorNetApy).Mul(seniorLiq)

		juniorApy := seniorLiq.Sub(a).Add(r.JuniorLiquidity).Mul(decimal.NewFromFloat(r.OriginatorNetApy)).Div(r.JuniorLiquidity)
		results[p.SmartYieldAddress].JuniorAPY, _ = juniorApy.Float64()
		results[p.SmartYieldAddress].AbondAPY = abondAPY
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_state", "included_in_block", "block_timestamp", "pool_address", "senior_liquidity", "junior_liquidity", "jtoken_price", "abond_principal", "abond_gain", "abond_issued_at", "abond_matures_at", "abond_apy", "senior_apy", "junior_apy", "originator_apy", "originator_net_apy"))
	if err != nil {
		return err
	}

	for _, p := range state.Pools() {
		r := results[p.SmartYieldAddress]

		_, err = stmt.Exec(s.Preprocessed.BlockNumber, s.Preprocessed.BlockTimestamp, r.PoolAddress, r.TotalLiquidity.Sub(r.JuniorLiquidity), r.JuniorLiquidity, r.JTokenPrice, r.Abond.Principal.String(), r.Abond.Gain.String(), decimal.NewFromBigInt(r.Abond.IssuedAt, -18).IntPart(), decimal.NewFromBigInt(r.Abond.MaturesAt, -18).IntPart(), r.AbondAPY, r.SeniorAPY, r.JuniorAPY, r.OriginatorApy, r.OriginatorNetApy)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}
