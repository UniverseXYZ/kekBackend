package smartYield

import (
	"database/sql"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type ERC721Transfer struct {
	*types.Event

	TokenAddress string
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
		TokenAddress: utils.CleanUpHex(log.Address),
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
	stmt, err := tx.Prepare(pq.CopyIn("erc721_transfers", "tx_hash", "tx_index", "log_index", "token_address", "sender", "receiver", "token_id", "included_in_block", "block_timestamp"))
	if err != nil {
		return err
	}

	for _, t := range s.processed.ERC721Transfers {
		_, err = stmt.Exec(t.TransactionHash, t.TransactionIndex, t.LogIndex, t.TokenAddress, t.Sender, t.Receiver, t.TokenID.String(), s.processed.blockNumber, s.processed.blockTimestamp)
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

	return nil
}
