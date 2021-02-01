package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/contracts"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleAbrogationProposal(logs []web3types.Log, tx *sql.Tx) error {
	var abrogationProposals []AbrogationProposal

	for _, log := range logs {
		if utils.LogIsEvent(log, g.govAbi, "AbrogationProposalStarted") {
			var cp AbrogationProposal
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

			err = g.govAbi.UnpackIntoInterface(&cp, "AbrogationProposalStarted", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			ctr, err := contracts.NewGovernance(common.HexToAddress(g.config.GovernanceAddress), g.ethConn)
			if err != nil {
				return err
			}

			p, err := ctr.AbrogationProposals(nil, proposalID)
			if err != nil {
				return errors.Wrap(err, "could not get the proposals from contract")
			}

			cp.ProposalID = proposalID
			cp.CreateTime = g.Preprocessed.BlockTimestamp
			cp.BaseLog = *baseLog
			cp.Description = p.Description
			abrogationProposals = append(abrogationProposals, cp)
		}
	}

	if len(abrogationProposals) == 0 {
		log.WithField("handler", "abrogation proposal").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_abrogation_proposals", "proposal_id", "creator", "create_time", "description", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, cp := range abrogationProposals {
		_, err = stmt.Exec(cp.ProposalID.Int64(), cp.Caller.String(), cp.CreateTime, cp.Description, cp.TransactionHash, cp.TransactionIndex, cp.LogIndex, cp.LoggedBy, g.Preprocessed.BlockNumber)
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
