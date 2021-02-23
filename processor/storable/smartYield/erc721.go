package smartYield

import (
	"database/sql"
	"math/big"
	"strings"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type ERC721Transfer struct {
	*types.Event

	SYAddress    string
	ProtocolId   string
	TokenAddress string
	TokenType    string
	Sender       string
	Receiver     string
	TokenID      *big.Int
}

func (s *Storable) decodeERC721TransferEvent(log web3types.Log) (*ERC721Transfer, error) {
	tokenID, ok := new(big.Int).SetString(utils.Trim0x(log.Topics[3]), 16)
	if !ok {
		return nil, errors.New("could not convert tokenID to big.Int")
	}

	event, err := new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &ERC721Transfer{
		TokenAddress: utils.NormalizeAddress(log.Address),
		Sender:       utils.Topic2Address(log.Topics[1]),
		Receiver:     utils.Topic2Address(log.Topics[2]),
		TokenID:      tokenID,
		Event:        event,
	}, nil
}

func (s *Storable) storeERC721Transfers(tx *sql.Tx) error {
	if len(s.processed.ERC721Transfers) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("erc721_transfers", "tx_hash", "tx_index", "log_index", "token_address", "sender", "receiver", "token_id", "token_type", "included_in_block", "block_timestamp"))
	if err != nil {
		return err
	}

	for _, t := range s.processed.ERC721Transfers {
		_, err = stmt.Exec(t.TransactionHash, t.TransactionIndex, t.LogIndex, t.TokenAddress, t.Sender, t.Receiver, t.TokenID.String(), t.TokenType, s.processed.blockNumber, s.processed.blockTimestamp)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	for _, a := range s.processed.ERC721Transfers {
		// we don't want to store Mint & Burn events in the transaction history table because there will already be another
		// corresponding action (eg: SeniorDeposit)
		if a.Sender == ZeroAddress || a.Receiver == ZeroAddress {
			continue
		}

		var tokenActionTypeSend, tokenActionTypeReceive string
		var amount decimal.Decimal
		if a.TokenType == "junior" {
			tokenActionTypeSend = string(JbondSend)
			tokenActionTypeReceive = string(JbondReceive)

			err := tx.QueryRow(`select tokens_in from smart_yield_junior_buy where junior_bond_id = $1`, a.TokenID.Int64()).Scan(&amount)
			if err != nil {
				return errors.Wrap(err, "could not find JuniorBond by id in the database")
			}
		} else {
			tokenActionTypeSend = string(SbondSend)
			tokenActionTypeReceive = string(SbondReceive)

			err := tx.QueryRow(`select underlying_in from smart_yield_senior_buy where senior_bond_id = $1`, a.TokenID.Int64()).Scan(&amount)
			if err != nil {
				return errors.Wrap(err, "could not find SeniorBond by id in the database")
			}
		}

		_, err = tx.Exec(`
			insert into smart_yield_transaction_history (protocol_id, sy_address, underlying_token_address, user_address, amount,
														 tranche, transaction_type, tx_hash, tx_index, log_index, block_timestamp,
														 included_in_block)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, a.ProtocolId, a.SYAddress, a.TokenAddress, a.Sender, amount, strings.ToUpper(a.TokenType), tokenActionTypeSend, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			insert into smart_yield_transaction_history (protocol_id, sy_address, underlying_token_address, user_address, amount,
														 tranche, transaction_type, tx_hash, tx_index, log_index, block_timestamp,
														 included_in_block)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, a.ProtocolId, a.SYAddress, a.TokenAddress, a.Receiver, amount, strings.ToUpper(a.TokenType), tokenActionTypeReceive, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	return nil
}
