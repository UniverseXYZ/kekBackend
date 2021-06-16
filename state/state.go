package state

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/types"
)

type State struct {
	db *sql.DB

	rewardPools       []types.SYRewardPool
	monitoredAccounts []string
	monitoredNTFs     []string
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
	err := loadAllAccounts()
	if err != nil {
		return errors.Wrap(err, "could not load monitored accounts ")
	}

	err = loadAllNFTs()
	if err != nil {
		return errors.Wrap(err, "could not load monitored nfts ")
	}

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

func loadAllNFTs() error {
	rows, err := instance.db.Query(`select address from monitored_nfts`)
	if err != nil {
		return errors.Wrap(err, "could not query database for monitored nfts")
	}

	var nfts []string
	for rows.Next() {
		var a string
		err := rows.Scan(&a)
		if err != nil {
			return errors.Wrap(err, "could no scan monitored nfts from database")
		}

		nfts = append(nfts, a)
	}

	instance.monitoredNTFs = nfts

	return nil
}
