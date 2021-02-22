package smartYieldState

import (
	"sync"

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
		results[p.SmartYieldAddress].TotalLiquidity = underlyingTotal
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
		results[p.SmartYieldAddress].JuniorLiquidity = underlyingJuniors
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
		results[p.SmartYieldAddress].JTokenPrice = price
		mu.Unlock()

		return nil
	})
}
