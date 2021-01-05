package processor

import (
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/metrics"
	"github.com/barnbridge/barnbridge-backend/processor/storable"
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
	"github.com/barnbridge/barnbridge-backend/processor/storable/governance"
	"github.com/barnbridge/barnbridge-backend/types"
)

var log = logrus.WithField("module", "data")

type Processor struct {
	config Config

	Raw *types.RawData

	abis    map[string]abi.ABI
	ethConn *ethclient.Client

	storables []Storable
}

func New(config Config, raw *types.RawData, abis map[string]abi.ABI, ethConn *ethclient.Client) (*Processor, error) {
	p := &Processor{
		config:  config,
		Raw:     raw,
		abis:    abis,
		ethConn: ethConn,
	}

	err := p.registerStorables()
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
func (fb *Processor) registerStorables() error {
	fb.storables = append(fb.storables, storable.NewStorableBlock(fb.Raw.Block))

	{
		if _, exist := fb.abis["bond"]; !exist {
			return errors.New("could not find abi for bond contract")
		}

		fb.storables = append(fb.storables, bond.NewBondStorable(fb.config.Bond, fb.Raw, fb.abis["bond"]))
	}

	{
		if _, exist := fb.abis["barn"]; !exist {
			return errors.New("could not find abi for barn contract")
		}
		fb.storables = append(fb.storables, barn.NewBarnStorable(fb.config.Barn, fb.Raw, fb.abis["barn"]))
	}

	{
		if _, exist := fb.abis["governance"]; !exist {
			return errors.New("could not find abi for governance contract")
		}
		fb.storables = append(fb.storables, governance.NewGovernance(fb.config.Governance, fb.Raw, fb.abis["governance"], fb.ethConn))
	}

	return nil
}

// Store will open a database transaction and execute all the registered Storables in the said transaction
func (fb *Processor) Store(db *sql.DB, m *metrics.Provider) error {
	exists, err := fb.checkBlockExists(db)
	if err != nil {
		return err
	}

	if exists {
		log.Info("block already exists in the database; skipping")
		return nil
	}

	reorged, err := fb.checkBlockReorged(db)
	if err != nil {
		return err
	}

	if reorged {
		m.RecordReorgedBlock()
		number, err := fb.extractBlockNumber()
		if err != nil {
			return err
		}
		log.WithField("block", number).Warn("detected reorged block")
		_, err = db.Exec("select delete_block($1)", number)
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

	for _, s := range fb.storables {
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
