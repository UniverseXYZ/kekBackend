package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleCancellationProposalVotes(logs []web3types.Log, tx *sql.Tx) error {
	var cancellationProposalVotes []Vote
	var cancellationProposalCancelledVotes []VoteCanceled
	for _, log := range logs {
		if utils.LogIsEvent(log, g.govAbi, "CancellationProposalVote") {
			var vote Vote
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			user := utils.Topic2Address(log.Topics[2])

			data, err := hex.DecodeString(utils.Trim0x(log.Data))
			if err != nil {
				return errors.Wrap(err, "could not decode log data")
			}

			err = g.govAbi.UnpackIntoInterface(&vote, "CancellationProposalVote", data)
			if err != nil {
				return errors.Wrap(err, "could not unpack log data")
			}

			baseLog, err := g.getBaseLog(log)
			if err != nil {
				return err
			}

			vote.ProposalID = proposalID
			vote.User = user
			vote.BaseLog = *baseLog
			vote.Timestamp = g.Preprocessed.BlockTimestamp
			cancellationProposalVotes = append(cancellationProposalVotes, vote)
		}

		if utils.LogIsEvent(log, g.govAbi, "CancellationProposalVoteCancelled") {
			var vote VoteCanceled
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}
			user := utils.Topic2Address(log.Topics[2])

			baseLog, err := g.getBaseLog(log)
			if err != nil {
				return err
			}

			vote.ProposalID = proposalID
			vote.User = user
			vote.BaseLog = *baseLog
			vote.Timestamp = g.Preprocessed.BlockTimestamp

			cancellationProposalCancelledVotes = append(cancellationProposalCancelledVotes, vote)
		}

	}

	err := g.insertCancellationProposalVotesToDB(cancellationProposalVotes, tx)
	if err != nil {
		return err
	}

	err = g.insertCancellationProposalVotesCanceledToDB(cancellationProposalCancelledVotes, tx)
	if err != nil {
		return err
	}

	return nil
}

func (g *GovStorable) insertCancellationProposalVotesToDB(votes []Vote, tx *sql.Tx) error {
	if len(votes) == 0 {
		log.WithField("handler", "cancellation proposal vote").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_cancellation_votes", "proposal_id", "user_id", "support", "power", "block_timestamp", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, v := range votes {
		_, err = stmt.Exec(v.ProposalID.Int64(), v.User, *v.Support, v.Power.String(), v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber)
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

func (g GovStorable) insertCancellationProposalVotesCanceledToDB(votes []VoteCanceled, tx *sql.Tx) error {
	if len(votes) == 0 {
		log.WithField("handler", "cancellation proposal cancel vote").Debug("no events found")
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("governance_cancellation_votes_canceled", "proposal_id", "user_id", "block_timestamp", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, v := range votes {
		_, err = stmt.Exec(v.ProposalID.Int64(), v.User, v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber)
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
