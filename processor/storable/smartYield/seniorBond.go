package smartYield

import (
	"database/sql"
	"encoding/hex"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type SeniorBondBuyTrade struct {
	*types.Event

	SYAddress              string
	ProtocolId             string
	UnderlyingTokenAddress string
	BuyerAddress           string
	SeniorBondID           *big.Int
	UnderlyingIn           *big.Int
	Gain                   *big.Int
	ForDays                *big.Int
}

type SeniorBondRedeemTrade struct {
	*types.Event

	SYAddress              string
	ProtocolId             string
	UnderlyingTokenAddress string
	OwnerAddress           string
	SeniorBondID           *big.Int
	Fee                    *big.Int
}

func (s *Storable) decodeSeniorBondBuyEvent(log web3types.Log, event string) (*SeniorBondBuyTrade, error) {
	pool := state.PoolBySmartYieldAddress(log.Address)

	var t SeniorBondBuyTrade
	t.SYAddress = pool.SmartYieldAddress
	t.UnderlyingTokenAddress = pool.UnderlyingAddress
	t.ProtocolId = pool.ProtocolId

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
		return nil, errors.New("could not convert seniorBondId ")
	}

	t.SeniorBondID = n
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) decodeSeniorBondRedeemEvent(log web3types.Log, event string) (*SeniorBondRedeemTrade, error) {
	pool := state.PoolBySmartYieldAddress(log.Address)

	var t SeniorBondRedeemTrade
	t.SYAddress = pool.SmartYieldAddress
	t.UnderlyingTokenAddress = pool.UnderlyingAddress
	t.ProtocolId = pool.ProtocolId

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
		return nil, errors.New("could not convert seniorBondId ")
	}

	t.SeniorBondID = n
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) storeSeniorBuyTrades(tx *sql.Tx) error {
	if len(s.processed.seniorActions.seniorBondBuys) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_senior_buy", "sy_address", "buyer_address", "senior_bond_address", "senior_bond_id", "underlying_in", "gain", "for_days", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.seniorActions.seniorBondBuys {
		_, err = stmt.Exec(a.SYAddress, a.BuyerAddress, state.PoolBySmartYieldAddress(a.SYAddress).SeniorBondAddress, a.SeniorBondID.String(), a.UnderlyingIn.String(), a.Gain.String(), a.ForDays.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

	for _, a := range s.processed.seniorActions.seniorBondBuys {
		_, err = tx.Exec(`
			insert into smart_yield_transaction_history (protocol_id, sy_address, underlying_token_address, user_address, amount,
														 tranche, transaction_type, tx_hash, tx_index, log_index, block_timestamp,
														 included_in_block)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, a.ProtocolId, a.SYAddress, a.UnderlyingTokenAddress, a.BuyerAddress, a.UnderlyingIn.String(), "SENIOR", SeniorDeposit, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storable) storeSeniorRedeemTrades(tx *sql.Tx) error {
	if len(s.processed.seniorActions.seniorBondRedeems) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_senior_redeem", "sy_address", "owner_address", "senior_bond_address", "senior_bond_id", "fee", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.seniorActions.seniorBondRedeems {
		_, err = stmt.Exec(a.SYAddress, a.OwnerAddress, state.PoolBySmartYieldAddress(a.SYAddress).SeniorBondAddress, a.SeniorBondID.String(), a.Fee.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

	for _, a := range s.processed.seniorActions.seniorBondRedeems {
		var underlyingIn, gain decimal.Decimal
		err := tx.QueryRow(`select underlying_in, gain from smart_yield_senior_buy where senior_bond_id = $1`, a.SeniorBondID.Int64()).Scan(&underlyingIn, &gain)
		if err != nil {
			return errors.Wrap(err, "could not find SeniorBond by id in the databas ")
		}

		amount := underlyingIn.Add(gain).Sub(decimal.NewFromBigInt(a.Fee, 0))

		_, err = tx.Exec(`
			insert into smart_yield_transaction_history (protocol_id, sy_address, underlying_token_address, user_address, amount,
														 tranche, transaction_type, tx_hash, tx_index, log_index, block_timestamp,
														 included_in_block)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`, a.ProtocolId, a.SYAddress, a.UnderlyingTokenAddress, a.OwnerAddress, amount, "SENIOR", SeniorRedeem, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	return nil
}
