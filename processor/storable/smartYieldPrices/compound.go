package smartYieldPrices

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) getCompoundPrice(oracleAddress string, wg *errgroup.Group, p types.SYPool, mu *sync.Mutex) {
	wg.Go(func() error {
		input, err := utils.ABIGenerateInput(s.abis["compoundoracle"], "getUnderlyingPrice", common.HexToAddress(p.ReceiptTokenAddress))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not generate input for %s.%s", oracleAddress, "getUnderlyingPrice"))
		}

		cp, err := s.callSimpleFunctionWithInput(s.abis["compoundoracle"], oracleAddress, "getUnderlyingPrice", input)
		if err != nil {
			return errors.Wrap(err, "could not call oracle.price('COMP')")
		}

		tokenPriceDec := decimal.NewFromBigInt(cp, -int32(18+18-p.UnderlyingDecimals))

		mu.Lock()
		s.Processed.Prices[p.ProtocolId+"/"+p.UnderlyingAddress] = Price{
			ProtocolId:   p.ProtocolId,
			TokenAddress: p.UnderlyingAddress,
			TokenSymbol:  p.UnderlyingSymbol,
			Value:        tokenPriceDec,
		}
		mu.Unlock()

		return nil
	})
}
