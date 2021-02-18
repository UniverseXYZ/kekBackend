package smartYieldState

import (
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type Config struct {
	NodeURL            string
	ComptrollerAddress string
}

type Storable struct {
	config Config
	raw    *types.RawData

	abis map[string]abi.ABI
	eth  *ethrpc.ETH

	Preprocessed struct {
		BlockTimestamp int64
		BlockNumber    int64
	}
}

func New(config Config, raw *types.RawData, abis map[string]abi.ABI) (*Storable, error) {
	var s = &Storable{
		config: config,
		raw:    raw,
		abis:   abis,
	}

	var err error
	s.Preprocessed.BlockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process block number")
	}

	s.Preprocessed.BlockTimestamp, err = strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse block timestamp")
	}

	batchLoader, err := httprpc.NewBatchLoader(0, 4*time.Millisecond)
	if err != nil {
		return nil, errors.Wrap(err, "could not init batch loader")
	}

	provider, err := httprpc.NewWithLoader(config.NodeURL, batchLoader)
	if err != nil {
		return nil, err
	}

	s.eth, err = ethrpc.New(provider)
	if err != nil {
		return nil, err
	}

	return s, nil
}

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

func (s Storable) ToDB(tx *sql.Tx) error {
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

		if p.ProtocolId == "compound/v2" {
			s.getCompoundAPY(wg, p, mu, results)
		}
	}

	err := wg.Wait()
	if err != nil {
		return err
	}

	spew.Dump(results)

	return nil
}

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

func (s Storable) getCompoundAPY(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex, results map[string]*State) {
	wg.Go(func() error {
		subWG := &errgroup.Group{}

		var compSpeeds *big.Int
		var totalBorrows *big.Int
		var cash *big.Int
		var oracleAddress string

		subWG.Go(func() error {
			rate, err := s.callSimpleFunction(s.abis["ctoken"], p.ReceiptTokenAddress, "supplyRatePerBlock")
			if err != nil {
				return err
			}

			r := decimal.NewFromBigInt(rate, -18)

			rf, _ := r.Float64()
			blocksPerDay := float64(4 * 60 * 24)

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

		compDollarsPerBlock := decimal.NewFromBigInt(compSpeeds, -18).Mul(compPriceDec)

		totalSupplied := new(big.Int).Add(cash, totalBorrows)
		tsDec := decimal.NewFromBigInt(totalSupplied, -int32(p.UnderlyingDecimals)).Mul(tokenPriceDec)

		apr := compDollarsPerBlock.DivRound(tsDec, 18)

		rf, _ := apr.Float64()
		blocksPerDay := float64(4 * 60 * 24)

		apy := math.Pow(rf*blocksPerDay+1, 365) - 1

		results[p.SmartYieldAddress].OriginatorNetApy = results[p.SmartYieldAddress].OriginatorApy + apy

		return nil
	})
}

func (s Storable) callSimpleFunction(a abi.ABI, contract string, name string) (*big.Int, error) {
	input, err := utils.ABIGenerateInput(a, name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", contract, name))
	}

	return s.callSimpleFunctionWithInput(a, contract, name, input)
}

func (s Storable) callSimpleFunctionWithInput(a abi.ABI, contract string, name string, input string) (*big.Int, error) {
	data, err := utils.CallAtBlock(s.eth, contract, input, s.Preprocessed.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not call %s.%s", contract, name))
	}

	decoded, err := utils.DecodeFunctionOutput(a, name, data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode output from %s.%s", contract, name))
	}

	return decoded[""].(*big.Int), nil
}
