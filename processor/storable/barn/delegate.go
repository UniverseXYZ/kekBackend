package barn

import (
	"context"
	"database/sql"
	"encoding/hex"
	"time"

	web3types "github.com/alethio/web3-go/types"
	"github.com/kekDAO/kekBackend/notifications"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/kekDAO/kekBackend/utils"
)

// Delegate followed by DelegateLockDecreased + DelegateLockIncreased => user had a delegate and moved it to another user
// Delegate followed by DelegateLockDecreased => user called stopDelegate
// Delegate followed by DelegateLockIncreased => user called delegate without previously delegating
// single DelegateLockIncreased => user deposited some more BOND which was automatically delegated
// single DelegateLockDecreased => user withdrew some BOND
func (b *BarnStorable) handleDelegate(logs []web3types.Log, tx *sql.Tx) error {
	var delegateActions []DelegateAction
	var delegateChanges []DelegateChange
	var jobs []*notifications.Job

	for _, l := range logs {
		action, err := b.decodeDelegateEvent(l)
		if err != nil {
			return err
		} else if action != nil {
			delegateActions = append(delegateActions, *action)
		}

		increase, err := b.decodeDelegatePowerIncreased(l)
		if err != nil {
			return err
		} else if increase != nil {
			delegateChanges = append(delegateChanges, *increase)
		}

		decrease, err := b.decodeDelegatePowerDecreased(l)
		if err != nil {
			return err
		} else if decrease != nil {
			delegateChanges = append(delegateChanges, *decrease)
		}

		if increase != nil {
			if increase.ToNewDelegatedPower.Cmp(increase.Amount) == 0 {
				jd := notifications.DelegateJobData{
					StartTime:             b.Preprocessed.BlockTimestamp,
					From:                  increase.Sender,
					To:                    increase.Receiver,
					Amount:                decimal.NewFromBigInt(increase.ToNewDelegatedPower, 0),
					IncludedInBlockNumber: b.Preprocessed.BlockNumber,
				}
				j, err := notifications.NewDelegateStartJob(&jd)
				if err != nil {
					return errors.Wrap(err, "could not create notification job")
				}

				jobs = append(jobs, j)
			}
		}
	}

	err := b.storeDelegateActions(delegateActions, tx)
	if err != nil {
		return err
	}

	err = b.storeDelegateChanges(delegateChanges, tx)
	if err != nil {
		return err
	}

	if b.config.Notifications && len(jobs) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err = notifications.ExecuteJobsWithTx(ctx, tx, jobs...)
		if err != nil && err != context.DeadlineExceeded {
			return errors.Wrap(err, "could not execute notification jobs")
		}
	}

	return nil
}

func (b *BarnStorable) storeDelegateActions(actions []DelegateAction, tx *sql.Tx) error {
	if len(actions) == 0 {
		log.WithField("handler", "delegate actions").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("barn_delegate_actions", "tx_hash", "tx_index", "log_index", "logged_by", "sender", "receiver", "action_type", "timestamp", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.LoggedBy, a.Sender, a.Receiver, a.ActionType, b.Preprocessed.BlockTimestamp, b.Preprocessed.BlockNumber)
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

func (b *BarnStorable) storeDelegateChanges(changes []DelegateChange, tx *sql.Tx) error {
	if len(changes) == 0 {
		log.WithField("handler", "delegate changes").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("barn_delegate_changes", "tx_hash", "tx_index", "log_index", "logged_by", "action_type", "sender", "receiver", "amount", "receiver_new_delegated_power", "timestamp", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, c := range changes {
		_, err = stmt.Exec(c.TransactionHash, c.TransactionIndex, c.LogIndex, c.LoggedBy, c.ActionType, c.Sender, c.Receiver, c.Amount.String(), c.ToNewDelegatedPower.String(), b.Preprocessed.BlockTimestamp, b.Preprocessed.BlockNumber)
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

func (b *BarnStorable) decodeDelegateEvent(log web3types.Log) (*DelegateAction, error) {
	if !utils.LogIsEvent(log, b.barnAbi, DelegateEvent) {
		return nil, nil
	}

	baseLog, err := b.getBaseLog(log)
	if err != nil {
		return nil, err
	}

	sender := utils.Topic2Address(log.Topics[1])
	receiver := utils.Topic2Address(log.Topics[2])

	var action ActionType
	if receiver == ZeroAddress {
		action = DELEGATE_STOP
	} else {
		action = DELEGATE_START
	}

	return &DelegateAction{
		BaseLog:    *baseLog,
		Sender:     sender,
		Receiver:   receiver,
		ActionType: action,
	}, nil

}

func (b *BarnStorable) decodeDelegatePowerIncreased(log web3types.Log) (*DelegateChange, error) {
	if !utils.LogIsEvent(log, b.barnAbi, DelegatePowerIncreasedEvent) {
		return nil, nil
	}

	baseLog, err := b.getBaseLog(log)
	if err != nil {
		return nil, err
	}

	d := &DelegateChange{
		BaseLog:    *baseLog,
		ActionType: DELEGATE_INCREASE,
		Sender:     utils.Topic2Address(log.Topics[1]),
		Receiver:   utils.Topic2Address(log.Topics[2]),
	}

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.barnAbi.UnpackIntoInterface(d, DelegatePowerIncreasedEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	return d, nil
}

func (b *BarnStorable) decodeDelegatePowerDecreased(log web3types.Log) (*DelegateChange, error) {
	if !utils.LogIsEvent(log, b.barnAbi, DelegatePowerDecreasedEvent) {
		return nil, nil
	}

	baseLog, err := b.getBaseLog(log)
	if err != nil {
		return nil, err
	}

	d := &DelegateChange{
		BaseLog:    *baseLog,
		ActionType: DELEGATE_DECREASE,
		Sender:     utils.Topic2Address(log.Topics[1]),
		Receiver:   utils.Topic2Address(log.Topics[2]),
	}

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = b.barnAbi.UnpackIntoInterface(d, DelegatePowerDecreasedEvent, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	return d, nil
}
