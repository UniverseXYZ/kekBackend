package governance

import (
	"database/sql"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/contracts"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleProposals(logs []web3types.Log, tx *sql.Tx) error {
	var proposals []Proposal
	var actions []Action

	for _, log := range logs {
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCreated"].ID.String()) {
			ctr, err := contracts.NewGovernance(common.HexToAddress(g.config.GovernanceAddress), &g.GovernanceClient)
			if err != nil {
				return err
			}

			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			p, err := ctr.Proposals(nil, proposalID)
			if err != nil {
				return errors.Wrap(err, "could not get the proposals from contract")
			}

			proposals = append(proposals, p)

			a, err := ctr.GetActions(nil, proposalID)
			if err != nil {
				return errors.Wrap(err, "could not get the actions from contract")
			}

			actions = append(actions, a)
		}
	}
	stmt, err := tx.Prepare(pq.CopyIn("governance_proposals", "proposal_ID", "proposer", "description", "title", "create_time", "start_time", "quorum", "eta", "for_votes", "against_votes", "canceled", "executed", "targets", "values", "signatures", "calldatas", "included_in_block", "created_at"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, p := range proposals {

		_, err = stmt.Exec(p.Id, p.Proposer, p.Description, p.Title, p.CreateTime, p.StartTime, p.Quorum, p.Eta, p.ForVotes, p.AgainstVotes, p.Canceled, p.Executed, g.Preprocessed.BlockNumber, g.Preprocessed.BlockTimestamp)
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
