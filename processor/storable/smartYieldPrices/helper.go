package smartYieldPrices

import (
	"database/sql"
	"fmt"
	"math/big"
	"sync"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/sync/errgroup"
)

func (s Storable) callSimpleFunction(a abi.ABI, contract string, name string) (*big.Int, error) {
	input, err := utils.ABIGenerateInput(a, name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", contract, name))
	}

	return s.callSimpleFunctionWithInput(a, contract, name, input)
}

func (s Storable) callSimpleFunctionWithInput(a abi.ABI, contract string, name string, input string) (*big.Int, error) {
	data, err := utils.CallAtBlock(s.eth, contract, input, s.Processed.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not call %s.%s", contract, name))
	}

	decoded, err := utils.DecodeFunctionOutput(a, name, data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode output from %s.%s", contract, name))
	}

	return decoded[""].(*big.Int), nil
}

func (s *Storable) getPrice(wg *errgroup.Group, p types.SYPool, mu *sync.Mutex) {
	wg.Go(func() error {
		price, err := s.callSimpleFunction(s.abis["smartyield"], p.SmartYieldAddress, "price")
		if err != nil {
			return err
		}

		mu.Lock()
		s.Processed.Prices[p.UnderlyingAddress] = price
		mu.Unlock()

		return nil
	})
}

func (s Storable) storePrices(tx *sql.Tx) error {
	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_prices", "token_address", "price_usd", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for key, a := range s.Processed.Prices {
		_, err = stmt.Exec(key, a.String(), s.Processed.BlockTimestamp, s.Processed.BlockNumber)
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
