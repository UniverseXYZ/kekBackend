package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleVotes(logs []web3types.Log, tx *sql.Tx) error {
	var votes []Vote
	var canceledVotes []VoteCanceled
	for _, log := range logs {
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["Vote"].ID.String()) {
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

			err = g.govAbi.UnpackIntoInterface(&vote, "vote", data)
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
			votes = append(votes, vote)
		}
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["VoteCanceled"].ID.String()) {
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

			canceledVotes = append(canceledVotes, vote)
		}

	}

	if len(votes) == 0 {
		log.Debug("no events found")
		return nil
	}

	err := g.insertVotesToDB(votes, tx)
	if err != nil {
		return err
	}

	if len(canceledVotes) == 0 {
		log.Debug("no events found")
		return nil
	}

	err = g.insertVotesCanceledToDB(canceledVotes, tx)
	if err != nil {
		return err
	}

	return nil
}

func (g *GovStorable) insertVotesToDB(votes []Vote, tx *sql.Tx) error {
	stmt, err := tx.Prepare(pq.CopyIn("governance_votes", "proposal_id", "user_id", "support", "power", "block_timestamp", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, v := range votes {
		_, err = stmt.Exec(v.ProposalID, v.User, v.Support, v.Power, v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber)
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

func (g GovStorable) insertVotesCanceledToDB(votes []VoteCanceled, tx *sql.Tx) error {
	stmt, err := tx.Prepare(pq.CopyIn("governance_votes_canceled", "proposal_id", "user_id", "block_timestamp", "tx_hash", "tx_index", "log_index", "logged_by", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "could not prepare statement")
	}

	for _, v := range votes {
		_, err = stmt.Exec(v.ProposalID, v.User, v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber)
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
