package state

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type State struct {
	db *sql.DB

	pools []types.SYPool
}

var instance *State

func Init(db *sql.DB) error {
	if instance != nil {
		return nil
	}

	instance = &State{db: db}

	err := loadAllSYPools()
	if err != nil {
		return errors.Wrap(err, "could not load SmartYield pools")
	}

	return nil
}

func loadAllSYPools() error {
	rows, err := instance.db.Query(`select protocol_id, controller_address, model_address, provider_address, sy_address, oracle_address, junior_bond_address, senior_bond_address, receipt_token_address, underlying_address, underlying_symbol, underlying_decimals from smart_yield_pools;`)
	if err != nil {
		return errors.Wrap(err, "could not query database for SmartYield pools")
	}

	var pools []types.SYPool
	for rows.Next() {
		var p types.SYPool
		err := rows.Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.ReceiptTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals)
		if err != nil {
			return errors.Wrap(err, "could not scan pools from database")
		}

		p.ControllerAddress = utils.NormalizeAddress(p.ControllerAddress)
		p.ModelAddress = utils.NormalizeAddress(p.ModelAddress)
		p.ProviderAddress = utils.NormalizeAddress(p.ProviderAddress)
		p.SmartYieldAddress = utils.NormalizeAddress(p.SmartYieldAddress)
		p.OracleAddress = utils.NormalizeAddress(p.OracleAddress)
		p.JuniorBondAddress = utils.NormalizeAddress(p.JuniorBondAddress)
		p.SeniorBondAddress = utils.NormalizeAddress(p.SeniorBondAddress)
		p.ReceiptTokenAddress = utils.NormalizeAddress(p.ReceiptTokenAddress)
		p.UnderlyingAddress = utils.NormalizeAddress(p.UnderlyingAddress)

		pools = append(pools, p)
	}

	instance.pools = pools

	return nil
}
