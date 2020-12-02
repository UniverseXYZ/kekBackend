package bond

import (
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/barnbridge/barnbridge-backend/types"
)

type BondStorable struct {
	config  Config
	Raw     *types.RawData
	bondAbi abi.ABI
}

func NewBondStorable(config Config, raw *types.RawData, bondAbi abi.ABI) *BondStorable {
	return &BondStorable{
		config:  config,
		Raw:     raw,
		bondAbi: bondAbi,
	}
}

func (b BondStorable) ToDB(tx *sql.Tx) error {
	fmt.Println(b.config.BondAddress)

	return nil
}
