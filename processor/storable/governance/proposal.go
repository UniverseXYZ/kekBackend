package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/barnbridge/barnbridge-backend/notifications"
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
			ctr, err := contracts.NewGovernance(common.HexToAddress(g.config.GovernanceAddress), g.ethConn)
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

			proposals = append(proposals, Proposal{
				Id:           p.Id,
				Proposer:     p.Proposer,
				Description:  p.Description,
				Title:        p.Title,
				CreateTime:   p.CreateTime,
				Eta:          p.Eta,
				ForVotes:     p.ForVotes,
				AgainstVotes: p.AgainstVotes,
				Canceled:     p.Canceled,
				Executed:     p.Executed,
				ProposalParameters: ProposalParameters{
					WarmUpDuration:      p.Parameters.WarmUpDuration,
					ActiveDuration:      p.Parameters.ActiveDuration,
					QueueDuration:       p.Parameters.QueueDuration,
					GracePeriodDuration: p.Parameters.GracePeriodDuration,
					AcceptanceThreshold: p.Parameters.AcceptanceThreshold,
					MinQuorum:           p.Parameters.MinQuorum,
				},
			})

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

	var jobs []*notifications.Job

	stmt, err := tx.Prepare(pq.CopyIn("governance_proposals", "proposal_id", "proposer", "description", "title", "create_time", "targets", "values", "signatures", "calldatas", "warm_up_duration", "active_duration", "queue_duration", "grace_period_duration", "acceptance_threshold", "min_quorum", "included_in_block", "block_timestamp"))
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

		_, err = stmt.Exec(p.Id.Int64(), p.Proposer.String(), p.Description, p.Title, p.CreateTime.Int64(), targets, values, signatures, calldatas, p.WarmUpDuration.Int64(), p.ActiveDuration.Int64(), p.QueueDuration.Int64(), p.GracePeriodDuration.Int64(), p.AcceptanceThreshold.Int64(), p.MinQuorum.Int64(), g.Preprocessed.BlockNumber, g.Preprocessed.BlockTimestamp)
		if err != nil {
			return errors.Wrap(err, "could not execute statement")
		}

		j, err := notifications.ProposalCreatedJob(p.Id.Int64())
		if err != nil {
			return errors.Wrap(err, "could not create notification job")
		}

		jobs = append(jobs, j)
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "could not close statement")
	}

	err = notifications.ExecuteJobsWithTx(tx, jobs...)
	if err != nil {
		return errors.Wrap(err, "could not execute notification jobs")
	}

	return nil
}
