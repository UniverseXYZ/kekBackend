package governance

import (
	"database/sql"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lib/pq"
	"github.com/pkg/errors"

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

			events = append(events, ProposalEvent{
				BaseLog:    *baseLog,
				ProposalID: proposalID,
				EventType:  CREATED,
			})

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalQueued"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, ProposalEvent{
				BaseLog:    *baseLog,
				ProposalID: proposalID,
				EventType:  QUEUED,
			})

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalExecuted"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, ProposalEvent{
				BaseLog:    *baseLog,
				ProposalID: proposalID,
				EventType:  EXECUTED,
			})

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCanceled"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, ProposalEvent{
				BaseLog:    *baseLog,
				ProposalID: proposalID,
				EventType:  CANCELED,
			})

			continue
		}

	}

	if len(events) == 0 {
		log.Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_events", "proposal_ID", "event_type", "timestamp", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, e := range events {

		_, err = stmt.Exec(e.ProposalID, e.EventType, g.Preprocessed.BlockTimestamp, e.TransactionHash, e.TransactionIndex, e.LogIndex, e.LoggedBy, g.Preprocessed.BlockNumber)
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
