package smartYieldState

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s Storable) getTotalLiquidity(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		underlyingTotal, err := s.callSimpleFunction(s.abis["smartyield"], p.SmartYieldAddress, "underlyingTotal")
		if err != nil {
			return err
		}

		mu.Lock()
		results[p.SmartYieldAddress].TotalLiquidity = decimal.NewFromBigInt(underlyingTotal, 0)
		mu.Unlock()

		return nil
	})
}

func (s Storable) getJuniorLiquidity(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		underlyingJuniors, err := s.callSimpleFunction(s.abis["smartyield"], p.SmartYieldAddress, "underlyingJuniors")
		if err != nil {
			return err
		}

		mu.Lock()
		results[p.SmartYieldAddress].JuniorLiquidity = decimal.NewFromBigInt(underlyingJuniors, 0)
		mu.Unlock()

		return nil
	})
}

func (s Storable) getPrice(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		price, err := s.callSimpleFunction(s.abis["smartyield"], p.SmartYieldAddress, "price")
		if err != nil {
			return err
		}

		mu.Lock()
		results[p.SmartYieldAddress].JTokenPrice = decimal.NewFromBigInt(price, 0)
		mu.Unlock()

		return nil
	})
}

func (s Storable) getMaxBondDailyRate(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		r, err := s.callSimpleFunction(s.abis["smartyield"], p.SmartYieldAddress, "maxBondDailyRate")
		if err != nil {
			if strings.Contains(err.Error(), "VM execution error Reverted") {
				mu.Lock()
				results[p.SmartYieldAddress].SeniorAPY = 0
				mu.Unlock()

				return nil
			}

			return err
		}

		rate, _ := decimal.NewFromBigInt(r, -18).Float64()

		apy := math.Pow(rate+1, 365) - 1

		mu.Lock()
		results[p.SmartYieldAddress].SeniorAPY = apy
		mu.Unlock()

		return nil
	})
}

func (s Storable) getAbond(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		input, err := utils.ABIGenerateInput(s.abis["smartyield"], "abond")
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", p.SmartYieldAddress, "abond"))
		}

		data, err := utils.CallAtBlock(s.eth, p.SmartYieldAddress, input, s.Preprocessed.BlockNumber)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not call %s.%s", p.SmartYieldAddress, "abond"))
		}

		decoded, err := utils.DecodeFunctionOutput(s.abis["smartyield"], "abond", data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not decode output from %s.%s", p.SmartYieldAddress, "abond"))
		}

		mu.Lock()
		results[p.SmartYieldAddress].Abond = Abond{
			Principal:  decoded["principal"].(*big.Int),
			Gain:       decoded["gain"].(*big.Int),
			MaturesAt:  decoded["maturesAt"].(*big.Int),
			IssuedAt:   decoded["issuedAt"].(*big.Int),
			Liquidated: decoded["liquidated"].(bool),
		}
		mu.Unlock()

		return nil
	})
}
