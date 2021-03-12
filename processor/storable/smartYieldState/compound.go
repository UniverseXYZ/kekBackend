package smartYieldState

import (
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s Storable) getCompoundAPY(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		subWG := &errgroup.Group{}

		var compSpeeds *big.Int
		var totalBorrows *big.Int
		var harvestCost *big.Int
		var cash *big.Int
		var oracleAddress string

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
			input, err := utils.ABIGenerateInput(s.abis["comptroller"], "compSpeeds", common.HexToAddress(p.ReceiptTokenAddress))
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", s.config.ComptrollerAddress, ""))
			}

			rate, err := s.callSimpleFunctionWithInput(s.abis["comptroller"], s.config.ComptrollerAddress, "compSpeeds", input)
			if err != nil {
				return err
			}

			compSpeeds = rate

			return nil
		})

		subWG.Go(func() error {
			tb, err := s.callSimpleFunction(s.abis["ctoken"], p.ReceiptTokenAddress, "totalBorrows")
			if err != nil {
				return errors.Wrap(err, "could not get totalBorrows")
			}

			totalBorrows = tb

			return nil
		})

		subWG.Go(func() error {
			c, err := s.callSimpleFunction(s.abis["ctoken"], p.ReceiptTokenAddress, "getCash")
			if err != nil {
				return errors.Wrap(err, "could not get cash")
			}

			cash = c

			return nil
		})

		subWG.Go(func() error {
			input, err := utils.ABIGenerateInput(s.abis["comptroller"], "oracle")
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", s.config.ComptrollerAddress, "oracle"))
			}

			data, err := utils.CallAtBlock(s.eth, s.config.ComptrollerAddress, input, s.Preprocessed.BlockNumber)
			if err != nil {
				return errors.Wrap(err, "could not call comptroller.oracle()")
			}

			oracleAddress = utils.Topic2Address(data)

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

		var compPrice *big.Int
		var tokenPrice *big.Int

		subWG.Go(func() error {
			input, err := utils.ABIGenerateInput(s.abis["compoundoracle"], "price", "COMP")
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", oracleAddress, "price"))
			}

			cp, err := s.callSimpleFunctionWithInput(s.abis["compoundoracle"], oracleAddress, "price", input)
			if err != nil {
				return errors.Wrap(err, "could not call oracle.price('COMP')")
			}

			compPrice = cp

			return nil
		})

		subWG.Go(func() error {
			input, err := utils.ABIGenerateInput(s.abis["compoundoracle"], "getUnderlyingPrice", common.HexToAddress(p.ReceiptTokenAddress))
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", oracleAddress, "getUnderlyingPrice"))
			}

			cp, err := s.callSimpleFunctionWithInput(s.abis["compoundoracle"], oracleAddress, "getUnderlyingPrice", input)
			if err != nil {
				return errors.Wrap(err, "could not call oracle.price('COMP')")
			}

			tokenPrice = cp

			return nil
		})

		err = subWG.Wait()
		if err != nil {
			return errors.Wrap(err, "could not get compound distribution apy")
		}

		compPriceDec := decimal.NewFromBigInt(compPrice, -6)
		tokenPriceDec := decimal.NewFromBigInt(tokenPrice, -int32(18+18-p.UnderlyingDecimals))

		// subtract the HARVEST_COST from the COMP price
		harvestCostDec := decimal.NewFromBigInt(harvestCost, -18).Mul(compPriceDec)
		compPriceDec = compPriceDec.Sub(harvestCostDec)

		compDollarsPerBlock := decimal.NewFromBigInt(compSpeeds, -18).Mul(compPriceDec)

		totalSupplied := new(big.Int).Add(cash, totalBorrows)
		tsDec := decimal.NewFromBigInt(totalSupplied, -int32(p.UnderlyingDecimals)).Mul(tokenPriceDec)

		apr := compDollarsPerBlock.DivRound(tsDec, 18)

		rf, _ := apr.Float64()
		blocksPerDay := float64(s.config.BlocksPerMinute * 60 * 24)

		apy := math.Pow(rf*blocksPerDay+1, 365) - 1

		results[p.SmartYieldAddress].OriginatorNetApy = results[p.SmartYieldAddress].OriginatorApy + apy

		return nil
	})
}
