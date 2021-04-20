package processor

import (
	"database/sql"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kekDAO/kekBackend/state"

	"github.com/kekDAO/kekBackend/metrics"
	"github.com/kekDAO/kekBackend/processor/storable"
	"github.com/kekDAO/kekBackend/processor/storable/accountERC20Transfers"
	"github.com/kekDAO/kekBackend/processor/storable/governance"
	"github.com/kekDAO/kekBackend/processor/storable/kek"
	"github.com/kekDAO/kekBackend/processor/storable/supernova"
	"github.com/kekDAO/kekBackend/processor/storable/yieldFarming"
	"github.com/kekDAO/kekBackend/types"
)

var log = logrus.WithField("module", "data")

type Processor struct {
	config Config

	Raw *types.RawData

	abis     map[string]abi.ABI
	ethConn  *ethclient.Client
	ethBatch *ethrpc.ETH

	storables []Storable
}

func New(config Config, raw *types.RawData, abis map[string]abi.ABI, ethConn *ethclient.Client, ethBatch *ethrpc.ETH) (*Processor, error) {
	p := &Processor{
		config:   config,
		Raw:      raw,
		abis:     abis,
		ethConn:  ethConn,
		ethBatch: ethBatch,
	}

	err := state.Refresh()
	if err != nil {
		return nil, errors.Wrap(err, "could not refresh state")
	}

	err = p.registerStorables()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Storable
// role: a Storable serves as a means of transforming raw data and inserting it into the database
// input: raw Ethereum data + a database transaction
// output: processed/derived/enhanced data stored directly to the db
type Storable interface {
	ToDB(tx *sql.Tx) error
}

// RegisterStorables instantiates all the storables defined via code with the requested raw data
// Only the storables that are registered will be executed when the Store function is called
func (p *Processor) registerStorables() error {
	p.storables = append(p.storables, storable.NewStorableBlock(p.Raw.Block))

	{
		if _, exist := p.abis["kek"]; !exist {
			return errors.New("could not find abi for kek contract")
		}

		p.storables = append(p.storables, kek.NewKekStorable(p.config.Kek, p.Raw, p.abis["kek"]))
	}

	{
		if _, exist := p.abis["supernova"]; !exist {
			return errors.New("could not find abi for supernova contract")
		}
		p.storables = append(p.storables, supernova.NewSupernovaStorable(p.config.Supernova, p.Raw, p.abis["supernova"]))
	}

	{
		if _, exist := p.abis["governance"]; !exist {
			return errors.New("could not find abi for governance contract")
		}
		p.storables = append(p.storables, governance.NewGovernance(p.config.Governance, p.Raw, p.abis["governance"], p.ethConn))
	}

	{
		if _, exist := p.abis["yieldfarming"]; !exist {
			errors.New("could not find abi for yield farming contract")
		}
		p.storables = append(p.storables, yieldFarming.NewStorable(p.config.YieldFarming, p.Raw, p.abis["yieldfarming"]))
	}

	{
		if _, exists := p.abis["erc20"]; !exists {
			return errors.New("could not find erc20 abi")
		}

		p.storables = append(p.storables, accountERC20Transfers.NewStorable(p.config.AccountErc20Transfers, p.Raw, p.abis["erc20"], p.ethConn))
	}

	return nil
}

// Store will open a database transaction and execute all the registered Storables in the said transaction
func (p *Processor) Store(db *sql.DB, m *metrics.Provider) error {
	exists, err := p.checkBlockExists(db)
	if err != nil {
		return err
	}

	if exists {
		log.Info("block already exists in the database; skipping")
		return nil
	}

	reorged, err := p.checkBlockReorged(db)
	if err != nil {
		return err
	}

	if reorged {
		m.RecordReorgedBlock()
		number, err := p.extractBlockNumber()
		if err != nil {
			return err
		}
		log.WithField("block", number).Warn("detected reorged block")
		_, err = db.Exec("select delete_block($1, $2)", number, pq.Array(dbTables))
		if err != nil {
			log.Error(err)
			return err
		}
		log.WithField("block", number).Info("removed old version from the db; will be replaced with new version")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, s := range p.storables {
		err = s.ToDB(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
