package smartYield

import (
	"database/sql"
	"encoding/hex"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type JuniorBondBuyTrade struct {
	*types.Event

	SYAddress    string
	BuyerAddress string
	JuniorBondID *big.Int
	TokensIn     *big.Int
	MaturesAt    *big.Int
}

type JuniorBondRedeemTrade struct {
	*types.Event

	SYAddress     string
	OwnerAddress  string
	JuniorBondID  *big.Int
	UnderlyingOut *big.Int
}

func (s *Storable) decodeJuniorBondBuyEvent(log web3types.Log, event string) (*JuniorBondBuyTrade, error) {
	var t JuniorBondBuyTrade
	t.SYAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(&t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.BuyerAddress = utils.Topic2Address(log.Topics[1])

	n, ok := new(big.Int).SetString(utils.Trim0x(log.Topics[2]), 16)
	if !ok {
		return nil, errors.New("could not convert JuniorBondID ")
	}

	t.JuniorBondID = n
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) decodeJuniorBondRedeemEvent(log web3types.Log, event string) (*JuniorBondRedeemTrade, error) {
	var t JuniorBondRedeemTrade
	t.SYAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(&t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.OwnerAddress = utils.Topic2Address(log.Topics[1])

	n, ok := new(big.Int).SetString(utils.Trim0x(log.Topics[2]), 16)
	if !ok {
		return nil, errors.New("could not convert JuniorBondID ")
	}

	t.JuniorBondID = n
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) storeJuniorBuyTrades(tx *sql.Tx) error {
	if len(s.processed.juniorActions.juniorBondBuys) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_junior_buy", "sy_address", "buyer_address", "junior_bond_id", "tokens_in", "matures_at", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.juniorActions.juniorBondBuys {
		_, err = stmt.Exec(a.SYAddress, a.BuyerAddress, a.JuniorBondID.String(), a.TokensIn.String(), a.MaturesAt.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

func (s *Storable) storeJuniorRedeemTrades(tx *sql.Tx) error {
	if len(s.processed.juniorActions.juniorBondRedeems) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_junior_redeem", "sy_address", "owner_address", "junior_bond_id", "underlying_out", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.juniorActions.juniorBondRedeems {
		_, err = stmt.Exec(a.SYAddress, a.OwnerAddress, a.JuniorBondID.String(), a.UnderlyingOut.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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
