package state

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

type State struct {
	db *sql.DB

	syPools           []types.SYPool
	rewardPools       []types.SYRewardPool
	monitoredAccounts []string
}

var instance *State

func Init(db *sql.DB) error {
	if instance != nil {
		return nil
	}

	instance = &State{db: db}

	return Refresh()
}

func Refresh() error {
	err := loadAllSYPools()
	if err != nil {
		return errors.Wrap(err, "could not load SmartYield pools")
	}

	err = loadAllRewardPools()
	if err != nil {
		return errors.Wrap(err, "could not load SmartYield Reward pools")
	}

	err = loadAllAccounts()
	if err != nil {
		return errors.Wrap(err, "could not load monitored accounts ")
	}

	return nil
}

func loadAllSYPools() error {
	rows, err := instance.db.Query(`select protocol_id, controller_address, model_address, provider_address, sy_address, oracle_address, junior_bond_address, senior_bond_address, receipt_token_address, underlying_address, underlying_symbol, underlying_decimals, start_at_block from smart_yield_pools;`)
	if err != nil {
		return errors.Wrap(err, "could not query database for SmartYield pools")
	}

	var pools []types.SYPool
	for rows.Next() {
		var p types.SYPool
		err := rows.Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.ReceiptTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals, &p.StartAtBlock)
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

	instance.syPools = pools

	return nil
}

func loadAllRewardPools() error {
	rows, err := instance.db.Query(`select pool_address,pool_token_address,reward_token_address, start_at_block from smart_yield_reward_pools;`)
	if err != nil {
		return errors.Wrap(err, "could not query database for SmartYield Reward pools")
	}

	var pools []types.SYRewardPool
	for rows.Next() {
		var p types.SYRewardPool
		err := rows.Scan(&p.PoolAddress, &p.PoolTokenAddress, &p.RewardTokenAddress, &p.StartAtBlock)
		if err != nil {
			return errors.Wrap(err, "could not scan reward pools from database")
		}
		p.PoolAddress = utils.NormalizeAddress(p.PoolAddress)
		p.PoolTokenAddress = utils.NormalizeAddress(p.PoolTokenAddress)
		p.RewardTokenAddress = utils.NormalizeAddress(p.RewardTokenAddress)

		pools = append(pools, p)
	}

	instance.rewardPools = pools

	return nil
}

func loadAllAccounts() error {
	rows, err := instance.db.Query(`select address from monitored_accounts`)
	if err != nil {
		return errors.Wrap(err, "could not query database for monitored accounts")
	}

	var accounts []string
	for rows.Next() {
		var a string
		err := rows.Scan(&a)
		if err != nil {
			return errors.Wrap(err, "could no scan monitored accounts from database")
		}

		accounts = append(accounts, a)
	}

	instance.monitoredAccounts = accounts

	return nil
}
