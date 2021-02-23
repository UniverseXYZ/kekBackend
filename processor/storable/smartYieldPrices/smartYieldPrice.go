package smartYieldPrices

import (
	"database/sql"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
)

type Config struct {
	NodeURL string
}

type Storable struct {
	config Config
	raw    *types.RawData

	abis map[string]abi.ABI
	eth  *ethrpc.ETH

	Processed struct {
		Prices         map[string]*big.Int
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
	s.Processed.BlockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process block number")
	}

	s.Processed.BlockTimestamp, err = strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse block timestamp")
	}

	s.Processed.Prices = make(map[string]*big.Int)

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

func (s Storable) ToDB(tx *sql.Tx) error {
	var wg = &errgroup.Group{}
	var mu = &sync.Mutex{}

	for _, p := range state.Pools() {
		s.getPrice(wg, p, mu)
	}

	err := wg.Wait()
	if err != nil {
		return err
	}

	err = s.storePrices(tx)
	if err != nil {
		return err
	}

	return nil
}
