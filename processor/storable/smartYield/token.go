package smartYield

import (
	"context"
	"database/sql"
	"encoding/hex"
	"math/big"
	"time"

	web3types "github.com/alethio/web3-go/types"
	"github.com/barnbridge/barnbridge-backend/notifications"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type TokenBuyTrade struct {
	*types.Event

	SYAddress              string
	ProtocolId             string
	UnderlyingTokenAddress string
	BuyerAddress           string
	UnderlyingIn           *big.Int
	TokensOut              *big.Int
	Fee                    *big.Int
}

type TokenSellTrade struct {
	*types.Event

	SYAddress              string
	ProtocolId             string
	UnderlyingTokenAddress string
	SellerAddress          string
	TokensIn               *big.Int
	UnderlyingOut          *big.Int
	Forfeits               *big.Int
}

func (s *Storable) decodeTokenBuyEvent(log web3types.Log, event string) (*TokenBuyTrade, error) {
	pool := state.PoolBySmartYieldAddress(log.Address)

	var t TokenBuyTrade
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
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) decodeTokenSellEvent(log web3types.Log, event string) (*TokenSellTrade, error) {
	pool := state.PoolBySmartYieldAddress(log.Address)

	var t TokenSellTrade
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

	t.SellerAddress = utils.Topic2Address(log.Topics[1])
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) storeTokenBuyTrades(tx *sql.Tx) error {
	if len(s.processed.tokenActions.tokenBuyTrades) == 0 {
		return nil
	}

	var jobs []*notifications.Job

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_token_buy", "sy_address", "buyer_address", "underlying_in", "tokens_out", "fee", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.tokenActions.tokenBuyTrades {
		_, err = stmt.Exec(a.SYAddress, a.BuyerAddress, a.UnderlyingIn.String(), a.TokensOut.String(), a.Fee.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}

		jd := notifications.SmartYieldJobData{
			StartTime:             s.processed.blockTimestamp,
			PoolAddress:           a.SYAddress,
			Buyer:                 a.BuyerAddress,
			Amount:                decimal.NewFromBigInt(a.TokensOut, 0),
			IncludedInBlockNumber: s.processed.blockNumber,
		}
		j, err := notifications.NewSmartYieldTokenBoughtJob(&jd)
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
		return err
	}

	for _, a := range s.processed.tokenActions.tokenBuyTrades {
		_, err = tx.Exec(`
			insert into smart_yield_transaction_history (protocol_id, sy_address, underlying_token_address, user_address, amount,
														 tranche, transaction_type, tx_hash, tx_index, log_index, block_timestamp,
														 included_in_block)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`, a.ProtocolId, a.SYAddress, a.UnderlyingTokenAddress, a.BuyerAddress, a.UnderlyingIn.String(), "JUNIOR", JuniorDeposit, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	if s.config.Notifications && len(jobs) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err = notifications.ExecuteJobsWithTx(ctx, tx, jobs...)
		if err != nil && err != context.DeadlineExceeded {
			return errors.Wrap(err, "could not execute notification jobs")
		}
	}

	return nil
}

func (s *Storable) storeTokenSellTrades(tx *sql.Tx) error {
	if len(s.processed.tokenActions.tokenSellTrades) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_token_sell", "sy_address", "seller_address", "tokens_in", "underlying_out", "forfeits", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.tokenActions.tokenSellTrades {
		_, err = stmt.Exec(a.SYAddress, a.SellerAddress, a.TokensIn.String(), a.UnderlyingOut.String(), a.Forfeits.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
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

	for _, a := range s.processed.tokenActions.tokenSellTrades {
		_, err = tx.Exec(`
			insert into smart_yield_transaction_history (protocol_id, sy_address, underlying_token_address, user_address, amount,
														 tranche, transaction_type, tx_hash, tx_index, log_index, block_timestamp,
														 included_in_block)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`, a.ProtocolId, a.SYAddress, a.UnderlyingTokenAddress, a.SellerAddress, a.TokensIn.String(), "JUNIOR", JuniorInstantWithdraw, a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	return nil
}
