package governance

import (
	"database/sql"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/barnbridge/barnbridge-backend/contracts"
	"github.com/barnbridge/barnbridge-backend/types"
)

type GovStorable struct {
	config           Config
	Raw              *types.RawData
	govAbi           abi.ABI
	BlockTimestamp   int64
	BlockNumber      int64
	GovernanceClient ethclient.Client
}

func NewGovernance(config Config, raw *types.RawData, govAbi abi.ABI) *GovStorable {
	return &GovStorable{
		config: config,
		Raw:    raw,
		govAbi: govAbi,
	}
}

func (g GovStorable) ToDB(tx *sql.Tx) error {
	ctr, err := contracts.NewGovernance(common.HexToAddress(g.config.GovernanceAddress), &g.GovernanceClient)
	if err != nil {
		return err
	}
	spew.Dump(ctr)
	//p := ctr.Proposals(nil, big.NewInt())
	return nil
}
