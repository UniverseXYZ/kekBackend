package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleCancellationProposals(logs []web3types.Log, tx *sql.Tx) error {
	var cancellationProposals []CancellationProposal

	for _, log := range logs {
		if utils.LogIsEvent(log, g.govAbi, "CancellationProposalStarted") {
			var cp CancellationProposal
			baseLog, err := g.getBaseLog(log)
			if err != nil {
				return err
			}

			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			data, err := hex.DecodeString(utils.Trim0x(log.Data))
			if err != nil {
				return errors.Wrap(err, "could not decode log data")
			}

			err = g.govAbi.UnpackIntoInterface(&cp, "CancellationProposalStarted", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			cp.ProposalID = *proposalID
			cp.CreateTime = g.Preprocessed.BlockTimestamp
			cp.BaseLog = *baseLog
			cancellationProposals = append(cancellationProposals, cp)
		}
	}

	if len(cancellationProposals) == 0 {
		log.Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_cancellation_proposals", "proposal_id", "creator", "create_time", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, cp := range cancellationProposals {

		_, err = stmt.Exec(cp.ProposalID, cp.Caller, cp.CreateTime, cp.TransactionHash, cp.TransactionIndex, cp.LogIndex, cp.LoggedBy, g.Preprocessed.BlockNumber)
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
