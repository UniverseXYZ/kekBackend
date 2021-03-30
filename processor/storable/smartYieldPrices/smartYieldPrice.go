package smartYieldPrices

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type Config struct {
	ComptrollerAddress string
}

type Price struct {
	ProtocolId   string
	TokenAddress string
	TokenSymbol  string
	Value        decimal.Decimal
}

type Storable struct {
	config Config
	raw    *types.RawData

	abis   map[string]abi.ABI
	eth    *ethrpc.ETH
	logger *logrus.Entry

	Processed struct {
		Prices         map[string]Price
		BlockTimestamp int64
		BlockNumber    int64
	}
}

func New(config Config, raw *types.RawData, abis map[string]abi.ABI, eth *ethrpc.ETH) (*Storable, error) {
	var s = &Storable{
		config: config,
		raw:    raw,
		abis:   abis,
		eth:    eth,
		logger: logrus.WithField("module", "storable(smartYieldPrice)"),
	}

	var err error
	s.Processed.BlockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process block number")
	}

	s.Processed.BlockTimestamp, err = strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse block timestamp")
	}

	s.Processed.Prices = make(map[string]Price)

	return s, nil
}

func (s Storable) ToDB(tx *sql.Tx) error {
	s.logger.Debug("executing")
	start := time.Now()
	defer func() {
		s.logger.WithField("duration", time.Since(start)).Debug("done")
	}()

	if len(state.Pools()) == 0 {
		return nil
	}

	var wg = &errgroup.Group{}
	var mu = &sync.Mutex{}

	minBlock := state.Pools()[0].StartAtBlock
	for _, p := range state.Pools() {
		if p.StartAtBlock < minBlock {
			minBlock = p.StartAtBlock
		}
	}

	if s.Processed.BlockNumber < minBlock {
		s.logger.Info("skipping block because there's no pool to scrape for")
		return nil
	}

	var compoundOracleAddress string
	wg.Go(func() error {
		input, err := utils.ABIGenerateInput(s.abis["comptroller"], "oracle")
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", s.config.ComptrollerAddress, "oracle"))
		}

		data, err := utils.CallAtBlock(s.eth, s.config.ComptrollerAddress, input, s.Processed.BlockNumber)
		if err != nil {
			return errors.Wrap(err, "could not call comptroller.oracle()")
		}

		compoundOracleAddress = utils.Topic2Address(data)

		return nil
	})

	err := wg.Wait()
	if err != nil {
		return err
	}

	for _, p := range state.Pools() {
		if s.Processed.BlockNumber < p.StartAtBlock {
			s.logger.WithField("pool", p.SmartYieldAddress).Info("skipping pool due to StartAtBlock property")
			continue
		}

		if p.ProtocolId == "compound/v2" {
			s.getCompoundPrice(compoundOracleAddress, wg, p, mu)
		}
	}

	err = wg.Wait()
	if err != nil {
		return err
	}

	err = s.storePrices(tx)
	if err != nil {
		return err
	}

	return nil
}
