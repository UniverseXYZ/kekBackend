package governance

import (
	"database/sql"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/contracts"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleCancellationProposals(logs []web3types.Log, tx *sql.Tx) error {
	var cancellationProposals []CancellationProposal

	for _, log := range logs {
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["CancellationProposalStarted"].ID.String()) {
			ctr, err := contracts.NewGovernance(common.HexToAddress(g.config.GovernanceAddress), &g.GovernanceClient)
			if err != nil {
				return err
			}

			baseLog, err := g.getBaseLog(log)
			if err != nil {
				return err
			}

			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			p, err := ctr.CancellationProposals(nil, proposalID)
			if err != nil {
				return errors.Wrap(err, "could not get the proposals from contract")
			}
			var cp CancellationProposal

			cp.ProposalID = *proposalID
			cp.CreateTime = p.CreateTime.Int64()
			cp.Creator = p.Creator.String()
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

		_, err = stmt.Exec(cp.ProposalID, cp.Creator, cp.CreateTime, cp.TransactionHash, cp.TransactionIndex, cp.LogIndex, cp.LoggedBy, g.Preprocessed.BlockNumber)
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
