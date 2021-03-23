package smartYieldState

import (
	"math"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/types"
)

func (s Storable) getCompoundAPY(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		subWG := &errgroup.Group{}

		var harvestCost *big.Int

		subWG.Go(func() error {
			rate, err := s.callSimpleFunction(s.abis["ctoken"], p.ReceiptTokenAddress, "supplyRatePerBlock")
			if err != nil {
				return err
			}

			r := decimal.NewFromBigInt(rate, -18)

			rf, _ := r.Float64()
			blocksPerDay := float64(s.config.BlocksPerMinute * 60 * 24)

			apy := math.Pow(rf*blocksPerDay+1, 365) - 1

			mu.Lock()
			results[p.SmartYieldAddress].OriginatorApy = apy
			mu.Unlock()

			return nil
		})

		subWG.Go(func() error {
			rate, err := s.callSimpleFunction(s.abis["compoundcontroller"], p.ControllerAddress, "spotDailyDistributionRateProvider")
			if err != nil {
				return errors.Wrap(err, "could not call CompoundController.spotDailyDistributionRateProvider")
			}

			r := decimal.NewFromBigInt(rate, -18)
			rf, _ := r.Float64()

			apy := math.Pow(rf+1, 365) - 1
			mu.Lock()
			results[p.SmartYieldAddress].OriginatorNetApy = apy
			mu.Unlock()

			return nil
		})

		subWG.Go(func() error {
			hc, err := s.callSimpleFunction(s.abis["compoundcontroller"], p.ControllerAddress, "HARVEST_COST")
			if err != nil {
				return errors.Wrap(err, "could not call CompoundController.HARVEST_COST")
			}

			harvestCost = hc

			return nil
		})

		err := subWG.Wait()
		if err != nil {
			return errors.Wrap(err, "could not get compound distribution apy")
		}

		// subtract the HARVEST_COST from the COMP price
		harvestCostDec := decimal.NewFromBigInt(harvestCost, -18)
		harvestCostFloat, _ := harvestCostDec.Float64()

		mu.Lock()
		apy := results[p.SmartYieldAddress].OriginatorApy
		distributionAPY := results[p.SmartYieldAddress].OriginatorNetApy * (1 - harvestCostFloat)

		results[p.SmartYieldAddress].OriginatorNetApy = apy + distributionAPY
		mu.Unlock()

		return nil
	})
}
