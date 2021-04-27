package governance

import (
	"context"
	"database/sql"
	"encoding/hex"
	"time"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kekDAO/kekBackend/notifications"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/kekDAO/kekBackend/contracts"
	"github.com/kekDAO/kekBackend/utils"
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

	var jobs []*notifications.Job

	stmt, err := tx.Prepare(pq.CopyIn("governance_abrogation_proposals", "proposal_id", "creator", "create_time", "description", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, cp := range abrogationProposals {
		_, err = stmt.Exec(cp.ProposalID.Int64(), cp.Caller.String(), cp.CreateTime, cp.Description, cp.TransactionHash, cp.TransactionIndex, cp.LogIndex, cp.LoggedBy, g.Preprocessed.BlockNumber)
		if err != nil {
			return errors.Wrap(err, "could not execute statement")
		}

		jd := notifications.AbrogationProposalCreatedJobData{
			Id:                    cp.ProposalID.Int64(),
			Proposer:              cp.Caller.String(),
			CreateTime:            cp.CreateTime,
			IncludedInBlockNumber: g.Preprocessed.BlockNumber,
		}
		j, err := notifications.NewAbrogationProposalCreatedJob(&jd)
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

	if g.config.Notifications {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		err = notifications.ExecuteJobsWithTx(ctx, tx, jobs...)
		if err != nil && err != context.DeadlineExceeded {
			return errors.Wrap(err, "could not execute notification jobs")
		}
	}

	return nil
}
