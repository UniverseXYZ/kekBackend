package processor

import (
	"database/sql"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldPrices"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldState"
	"github.com/barnbridge/barnbridge-backend/state"

	"github.com/barnbridge/barnbridge-backend/metrics"
	"github.com/barnbridge/barnbridge-backend/processor/storable"
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
	"github.com/barnbridge/barnbridge-backend/processor/storable/governance"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYield"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldRewards"
	"github.com/barnbridge/barnbridge-backend/processor/storable/yieldFarming"
	"github.com/barnbridge/barnbridge-backend/types"
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
		if _, exist := p.abis["bond"]; !exist {
			return errors.New("could not find abi for bond contract")
		}

		p.storables = append(p.storables, bond.NewBondStorable(p.config.Bond, p.Raw, p.abis["bond"]))
	}

	{
		if _, exist := p.abis["barn"]; !exist {
			return errors.New("could not find abi for barn contract")
		}
		p.storables = append(p.storables, barn.NewBarnStorable(p.config.Barn, p.Raw, p.abis["barn"]))
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
		if _, exist := p.abis["smartyield"]; !exist {
			return errors.New("could not find smartYield abi")
		}

		if _, exist := p.abis["juniorbond"]; !exist {
			return errors.New("could not find juniorbond  abi")
		}

		if _, exist := p.abis["seniorbond"]; !exist {
			return errors.New("could not find seniorbond abi")
		}

		if _, exist := p.abis["compoundprovider"]; !exist {
			return errors.New("could not find compound provider abi")
		}

		p.storables = append(p.storables, smartYield.NewStorable(p.config.SmartYield, p.Raw, p.abis))

		syState, err := smartYieldState.New(p.config.SmartYieldState, p.Raw, p.abis, p.ethBatch)
		if err != nil {
			return errors.Wrap(err, "could not initialize SmartYieldState storable")
		}
		p.storables = append(p.storables, syState)

		syPrices, err := smartYieldPrices.New(p.config.SmartYieldPrice, p.Raw, p.abis, p.ethBatch)
		if err != nil {
			return errors.Wrap(err, "could not initialize SmartYieldPrice storable")
		}
		p.storables = append(p.storables, syPrices)
	}

	{
		if _, exist := p.abis["syreward"]; !exist {
			return errors.New("could not find smart yield rewards abi")
		}

		if _, exist := p.abis["poolfactory"]; !exist {
			return errors.New("could not find pool factory abi")
		}
		p.storables = append(p.storables, smartYieldRewards.NewStorable(p.config.SmartYieldRewards, p.Raw, p.abis["syreward"], p.abis["poolfactory"], p.ethConn))
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
