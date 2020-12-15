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
	var events []Event

	for _, log := range logs {
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCreated"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, Event{proposalID, CREATED})

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalQueued"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, Event{proposalID, QUEUED})

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalExecuted"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, Event{proposalID, EXECUTED})

			continue
		}

		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCanceled"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			events = append(events, Event{proposalID, CANCELED})

			continue
		}

	}

	if len(events) == 0 {
		log.Debug("Nothing to process...")
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("governance_events", "proposal_ID", "event_type", "included_in_block", "created_at"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, e := range events {

		_, err = stmt.Exec(e.ProposerID, e.EventType, g.Preprocessed.BlockNumber, g.Preprocessed.BlockTimestamp)
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
