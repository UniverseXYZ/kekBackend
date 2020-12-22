package governance

import (
	"database/sql"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

var log = logrus.WithField("module", "storable(governance)")

type GovStorable struct {
	config Config
	Raw    *types.RawData
	govAbi abi.ABI

	GovernanceClient ethclient.Client

	Preprocessed struct {
		BlockTimestamp int64
		BlockNumber    int64
	}
}

func NewGovernance(config Config, raw *types.RawData, govAbi abi.ABI) *GovStorable {
	return &GovStorable{
		config: config,
		Raw:    raw,
		govAbi: govAbi,
	}
}

func (g *GovStorable) preprocess() error {
	var err error

	g.Preprocessed.BlockNumber, err = strconv.ParseInt(g.Raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "unable to process block number")
	}

	g.Preprocessed.BlockTimestamp, err = strconv.ParseInt(g.Raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse block timestamp")
	}

	return nil
}

func (g GovStorable) ToDB(tx *sql.Tx) error {
	err := g.preprocess()
	if err != nil {
		return err
	}

	var govLogs []web3types.Log
	for _, data := range g.Raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) == utils.CleanUpHex(g.config.GovernanceAddress) {
				govLogs = append(govLogs, log)
			}
		}
	}

	if len(govLogs) == 0 {
		log.Debug("no events found")
		return nil
	}

	err = g.handleProposals(govLogs, tx)
	if err != nil {
		return err
	}

	err = g.handleEvents(govLogs, tx)
	if err != nil {
		return err
	}

	err = g.handleVotes(govLogs, tx)
	if err != nil {
		return err
	}

	err = g.handleCancellationProposals(govLogs, tx)
	if err != nil {
		return err
	}

	err = g.handleCancellationProposalVotes(govLogs, tx)
	if err != nil {
		return err
	}

	return nil
}
