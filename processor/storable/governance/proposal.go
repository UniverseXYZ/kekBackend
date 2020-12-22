package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/contracts"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleProposals(logs []web3types.Log, tx *sql.Tx) error {
	var proposals []Proposal
	var actions []ProposalActions

	for _, log := range logs {
		if utils.LogIsEvent(log, g.govAbi, "ProposalCreated") {
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

	if len(proposals) == 0 {
		log.WithField("handler", "proposals").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_proposals", "proposal_id", "proposer", "description", "title", "create_time", "targets", "values", "signatures", "calldatas", "included_in_block", "block_timestamp"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for i, p := range proposals {
		a := actions[i]
		var targets, values, signatures, calldatas types.JSONStringArray

		for i := 0; i < len(a.Targets); i++ {
			targets = append(targets, a.Targets[i].String())
			values = append(values, a.Values[i].String())
			signatures = append(signatures, a.Signatures[i])
			calldatas = append(calldatas, hex.EncodeToString(a.Calldatas[i]))
		}

		_, err = stmt.Exec(p.Id.Int64(), p.Proposer, p.Description, p.Title, p.CreateTime.Int64(), targets, values, signatures, calldatas, g.Preprocessed.BlockNumber, g.Preprocessed.BlockTimestamp)
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
