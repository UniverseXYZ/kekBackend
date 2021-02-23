package smartYieldState

import (
	"math"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/types"
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
