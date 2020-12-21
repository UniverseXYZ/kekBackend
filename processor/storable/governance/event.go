package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleEvents(logs []web3types.Log, tx *sql.Tx) error {
	var events []ProposalEvent

	for _, log := range logs {

		baseLog, err := g.getBaseLog(log)
		if err != nil {
			return err
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCreated"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			var e ProposalEvent
			e.BaseLog = *baseLog
			e.ProposalID = proposalID
			e.EventType = CREATED

			data, err := hex.DecodeString(utils.Trim0x(log.Data))
			if err != nil {
				return errors.Wrap(err, "could not decode log data")
			}

			err = g.govAbi.UnpackIntoInterface(&e, "ProposalCreated", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			events = append(events, e)

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalQueued"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			var e ProposalEvent
			e.BaseLog = *baseLog
			e.ProposalID = proposalID
			e.EventType = QUEUED

			data, err := hex.DecodeString(utils.Trim0x(log.Data))
			if err != nil {
				return errors.Wrap(err, "could not decode log data")
			}

			err = g.govAbi.UnpackIntoInterface(&e, "ProposalQueued", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			events = append(events, e)
			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalExecuted"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			var e ProposalEvent
			e.BaseLog = *baseLog
			e.ProposalID = proposalID
			e.EventType = EXECUTED

			data, err := hex.DecodeString(utils.Trim0x(log.Data))
			if err != nil {
				return errors.Wrap(err, "could not decode log data")
			}

			err = g.govAbi.UnpackIntoInterface(&e, "ProposalExecuted", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			events = append(events, e)
			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCanceled"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			var e ProposalEvent
			e.BaseLog = *baseLog
			e.ProposalID = proposalID
			e.EventType = CANCELED

			data, err := hex.DecodeString(utils.Trim0x(log.Data))
			if err != nil {
				return errors.Wrap(err, "could not decode log data")
			}

			err = g.govAbi.UnpackIntoInterface(&e, "ProposalExecuted", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			events = append(events, e)
			continue
		}

	}

	if len(events) == 0 {
		log.Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_events", "proposal_id", "caller", "event_type", "block_timestamp", "tx_hash", "tx_index", "log_index", "logged_by", "event_data", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, e := range events {
		var eventData types.JSONObject

		if e.Eta != nil {
			eventData = make(types.JSONObject)
			eventData["eta"] = e.Eta.Int64()
		}
		_, err = stmt.Exec(e.ProposalID, *e.Caller, e.EventType, g.Preprocessed.BlockTimestamp, e.TransactionHash, e.TransactionIndex, e.LogIndex, e.LoggedBy, eventData, g.Preprocessed.BlockNumber)
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
