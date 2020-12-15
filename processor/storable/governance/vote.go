package governance

import (
	"database/sql"
	"encoding/hex"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g *GovStorable) handleVotes(logs []web3types.Log, tx *sql.Tx) error {
	var votes []Vote
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
			vote.Canceled = false
			votes = append(votes, vote)
		}
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["VoteCanceled"].ID.String()) {
			var vote Vote
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
			vote.Canceled = true

			votes = append(votes, vote)
		}

	}

	if len(votes) == 0 {
		log.Debug("no events found")
		return nil
	}

	for _, v := range votes {
		var timestamp int64
		err := tx.QueryRow(`select "timestamp" from governance_votes where proposal_ID = $1 AND user_ID = $2`, v.ProposalID, v.User).Scan(&timestamp)

		if err != nil {
			if err == sql.ErrNoRows {
				sqlInsert := `INSERT INTO governance_votes(proposal_ID,user_ID,"timestamp",tx_hash,tx_index,log_index,logged_by,included_in_block,support,canceled,power)
 							  VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);`

				_, err2 := tx.Exec(sqlInsert, v.ProposalID, v.User, v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber, v.Support, v.Canceled, v.Power)
				if err2 != nil {
					return errors.Wrap(err2, "could not execute statement")
				}
			} else {
				return err
			}

		} else if timestamp < v.Timestamp {
			if v.Canceled == false {
				sqlUpdate := `
					UPDATE governance_votes
					SET support = $3 ,canceled = &4,timestamp = $5, tx_hash = $6 , tx_index = $7 , log_index = $8 , logged_by = $9 , included_in_block = $10
					WHERE proposal_ID = $1 AND user_ID = $2;`

				_, err = tx.Exec(sqlUpdate, v.ProposalID, v.User, v.Support, v.Canceled, v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber)
				if err != nil {
					return err
				}

			} else {
				sqlUpdateCanceled := `
					UPDATE governance_votes
					SET canceled = &3,timestamp = $4, tx_hash = $5 , tx_index = $6 , log_index = $7 , logged_by = $8 , included_in_block = $9
					WHERE proposal_ID = $1 AND user_ID = $2;`

				_, err = tx.Exec(sqlUpdateCanceled, v.ProposalID, v.User, v.Canceled, v.Timestamp, v.TransactionHash, v.TransactionIndex, v.LogIndex, v.LoggedBy, g.Preprocessed.BlockNumber)
				if err != nil {
					return err
				}

			}
		}
	}
	return nil
}
