package supernova

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/utils"
)

func (b *SupernovaStorable) handleLocks(logs []web3types.Log, tx *sql.Tx) error {
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
		log.WithField("handler", "locks").Debug("no events found")
		return nil
	}

	log.WithField("handler", "locks").WithField("count", len(locks)).Trace("found Lock events")

	stmt, err := tx.Prepare(pq.CopyIn("supernova_locks", "tx_hash", "tx_index", "log_index", "logged_by", "user_address", "locked_until", "locked_at", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, l := range locks {
		_, err = stmt.Exec(l.TransactionHash, l.TransactionIndex, l.LogIndex, l.LoggedBy, l.User, l.Timestamp.Int64(), b.Preprocessed.BlockTimestamp, b.Preprocessed.BlockNumber)
		if err != nil {
			return errors.Wrap(err, "could not execute statement")
		}
	}

	log.Trace("flushing data to db")

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

func (b *SupernovaStorable) decodeLockEvent(log web3types.Log) (*Lock, error) {
	if !utils.LogIsEvent(log, b.supernovaAbi, LockEvent) {
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

	err = b.supernovaAbi.UnpackIntoInterface(&lock, LockEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	return &lock, nil
}
