package barn

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (b *BarnStorable) handleLocks(logs []web3types.Log, tx *sql.Tx) error {
	var locks []Lock

	for _, log := range logs {
		l, err := b.decodeLockEvent(log)
		if err != nil {
			return err
		}

		if l != nil {
			locks = append(locks, *l)
		}
	}

	// if no lock was identified, fail fast to avoid extra processing
	if len(locks) == 0 {
		log.Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("barn_locks", "tx_hash", "tx_index", "log_index", "logged_by", "user_address", "locked_until", "locked_at", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, l := range locks {
		_, err = stmt.Exec(l.TransactionHash, l.TransactionIndex, l.LogIndex, l.LoggedBy, l.User, l.LockedUntil, b.Preprocessed.BlockTimestamp, b.Preprocessed.BlockNumber)
		if err != nil {
			return errors.Wrap(err, "could not execute statement")
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "could not close statement")
	}

	return nil
}

func (b *BarnStorable) decodeLockEvent(log web3types.Log) (*Lock, error) {
	if !utils.LogIsEvent(log, b.barnAbi, LockEvent) {
		return nil, nil
	}

	baseLog, err := b.getBaseLog(log)
	if err != nil {
		return nil, err
	}

	var lock = Lock{
		BaseLog: *baseLog,
		User:    utils.Topic2Address(log.Topics[1]),
	}

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.barnAbi.UnpackIntoInterface(&lock, LockEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	return &lock, nil
}
