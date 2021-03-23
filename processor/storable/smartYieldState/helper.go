package smartYieldState

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

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
