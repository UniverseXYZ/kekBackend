package auction

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

var logger = logrus.WithField("module", "storable(auction)")

type Storable struct {
	config     Config
	raw        *types.RawData
	auctionAbi abi.ABI
}

type AuctionEvent struct {
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
	data             []byte
}

type LogAuctionCreated struct {
	AuctionId         *big.Int       "json:\"auctionId\""
	AuctionOwner      common.Address "json:\"auctionOwner\""
	NumberOfSlots     *big.Int       "json:\"numberOfSlots\""
	StartTime         *big.Int       "json:\"startTime\""
	EndTime           *big.Int       "json:\"endTime\""
	ResetTimer        *big.Int       "json:\"resetTimer\""
	SupportsWhitelist bool           "json:\"supportsWhitelist\""
	Time              *big.Int       "json:\"time\""
}

type LogERC721Deposit struct {
	Depositor    common.Address "json:\"depositor\""
	TokenAddress common.Address "json:\"tokenAddress\""
	TokenId      *big.Int       "json:\"tokenId\""
	AuctionId    *big.Int       "json:\"auctionId\""
	SlotIndex    *big.Int       "json:\"slotIndex\""
	NftSlotIndex *big.Int       "json:\"nftSlotIndex\""
	Time         *big.Int       "json:\"time\""
}

type LogERC721Withdrawal struct {
	Depositor    common.Address "json:\"depositor\""
	TokenAddress common.Address "json:\"tokenAddress\""
	TokenId      *big.Int       "json:\"tokenId\""
	AuctionId    *big.Int       "json:\"auctionId\""
	SlotIndex    *big.Int       "json:\"slotIndex\""
	NftSlotIndex *big.Int       "json:\"nftSlotIndex\""
	Time         *big.Int       "json:\"time\""
}

type LogBidSubmitted struct {
	Sender     common.Address "json:\"sender\""
	AuctionId  *big.Int       "json:\"auctionId\""
	CurrentBid *big.Int       "json:\"currentBid\""
	TotalBid   *big.Int       "json:\"totalBid\""
	Time       *big.Int       "json:\"time\""
}

type LogBidWithdrawal struct {
	Recipient common.Address "json:\"recipient\""
	AuctionId *big.Int       "json:\"auctionId\""
	Amount    *big.Int       "json:\"amount\""
	Time      *big.Int       "json:\"time\""
}

type LogBidMatched struct {
	AuctionId        *big.Int       "json:\"auctionId\""
	SlotIndex        *big.Int       "json:\"slotIndex\""
	SlotReservePrice *big.Int       "json:\"slotReservePrice\""
	WinningBidAmount *big.Int       "json:\"winningBidAmount\""
	Winner           common.Address "json:\"winner\""
	Time             *big.Int       "json:\"time\""
}

type LogAuctionExtended struct {
	AuctionId *big.Int "json:\"auctionId\""
	EndTime   *big.Int "json:\"endTime\""
	Time      *big.Int "json:\"time\""
}

type LogAuctionCanceled struct {
	AuctionId *big.Int "json:\"auctionId\""
	Time      *big.Int "json:\"time\""
}

type LogAuctionRevenueWithdrawal struct {
	Recipient common.Address "json:\"recipient\""
	AuctionId *big.Int       "json:\"auctionId\""
	Smount    *big.Int       "json:\"amount\""
	Time      *big.Int       "json:\"time\""
}

type LogSlotRevenueCaptured struct {
	AuctionId *big.Int       "json:\"auctionId\""
	SlotIndex *big.Int       "json:\"slotIndex\""
	Amount    *big.Int       "json:\"amount\""
	BidToken  common.Address "json:\"bidToken\""
	Time      *big.Int       "json:\"time\""
}

type LogERC721RewardsClaim struct {
	Claimer   common.Address "json:\"claimer\""
	AuctionId *big.Int       "json:\"auctionId\""
	SlotIndex *big.Int       "json:\"slotIndex\""
	Time      *big.Int       "json:\"time\""
}

type LogRoyaltiesWithdrawal struct {
	Amount *big.Int       "json:\"amount\""
	To     common.Address "json:\"to\""
	Token  common.Address "json:\"token\""
	Time   *big.Int       "json:\"time\""
}

func NewStorable(config Config, raw *types.RawData, auctionAbi abi.ABI) *Storable {
	return &Storable{
		config:     config,
		raw:        raw,
		auctionAbi: auctionAbi,
	}
}

func (a *Storable) ToDB(tx *sql.Tx) error {
	var auctionEvents []AuctionEvent
	var erc721DepositEvents []AuctionEvent
	var auctionCanceledEvents []AuctionEvent
	var erc721WithdrawEvents []AuctionEvent
	var submittedBidEvents []AuctionEvent
	var withdrawnBidEvents []AuctionEvent
	var auctionExtendedEvents []AuctionEvent
	var revenueWithdrawEvents []AuctionEvent
	var slotRevenueEvents []AuctionEvent
	var erc721ClaimEvents []AuctionEvent
	var bidMatchedEvents []AuctionEvent

	for _, data := range a.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(a.config.Address) {
				continue
			}
			if len(log.Topics) == 0 {
				continue
			}

			if utils.LogIsEvent(log, a.auctionAbi, AuctionCreated) {
				d, err := a.decodeLog(log, AuctionCreated)
				if err != nil {
					return err
				}

				auctionEvents = append(auctionEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, ERC721Deposit) {
				d, err := a.decodeLog(log, ERC721Deposit)
				if err != nil {
					return err
				}

				erc721DepositEvents = append(erc721DepositEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, ERC721Withdrawal) {
				d, err := a.decodeLog(log, ERC721Withdrawal)
				if err != nil {
					return err
				}
				erc721WithdrawEvents = append(erc721WithdrawEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, AuctionCanceled) {
				d, err := a.decodeLog(log, AuctionCanceled)
				if err != nil {
					return err
				}

				auctionCanceledEvents = append(auctionCanceledEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, BidSubmitted) {
				d, err := a.decodeLog(log, BidSubmitted)
				if err != nil {
					return err
				}

				submittedBidEvents = append(submittedBidEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, BidWithdrawal) {
				d, err := a.decodeLog(log, BidWithdrawal)
				if err != nil {
					return err
				}

				withdrawnBidEvents = append(withdrawnBidEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, AuctionExtended) {
				d, err := a.decodeLog(log, AuctionExtended)
				if err != nil {
					return err
				}

				auctionExtendedEvents = append(auctionExtendedEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, BidMatched) {
				d, err := a.decodeLog(log, BidMatched)
				if err != nil {
					return err
				}

				bidMatchedEvents = append(bidMatchedEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, AuctionRevenueWithdrawal) {
				d, err := a.decodeLog(log, AuctionRevenueWithdrawal)
				if err != nil {
					return err
				}

				revenueWithdrawEvents = append(revenueWithdrawEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, SlotRevenueCaptured) {
				d, err := a.decodeLog(log, SlotRevenueCaptured)
				if err != nil {
					return err
				}

				slotRevenueEvents = append(slotRevenueEvents, *d)
			}

			if utils.LogIsEvent(log, a.auctionAbi, ERC721RewardsClaim) {
				d, err := a.decodeLog(log, ERC721RewardsClaim)
				if err != nil {
					return err
				}

				erc721ClaimEvents = append(erc721ClaimEvents, *d)
			}
		}
	}

	if len(auctionEvents) > 0 {
		err := a.storeActions(tx, auctionEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(erc721DepositEvents) > 0 {
		err := a.storeErc721DepositEvents(tx, erc721DepositEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(erc721WithdrawEvents) > 0 {
		err := a.storeErc721WithdrawEvents(tx, erc721WithdrawEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(auctionCanceledEvents) > 0 {
		err := a.storeAuctionCanceledEvents(tx, auctionCanceledEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(submittedBidEvents) > 0 {
		err := a.storeSubmittedBidEvents(tx, submittedBidEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(withdrawnBidEvents) > 0 {
		err := a.storeWithdrawnBidEvents(tx, withdrawnBidEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(auctionExtendedEvents) > 0 {
		err := a.storeAuctionExtendedEvents(tx, auctionExtendedEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(bidMatchedEvents) > 0 {
		err := a.storeBidMatchedEvents(tx, bidMatchedEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(revenueWithdrawEvents) > 0 {
		err := a.storeRevenueWithdrawEvents(tx, revenueWithdrawEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	if len(slotRevenueEvents) > 0 {
		err := a.storeCaptureSlotRevenueEvents(tx, slotRevenueEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}


	if len(erc721ClaimEvents) > 0 {
		err := a.storeErc721ClaimEvents(tx, erc721ClaimEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "auction events").Debug("no event found")
	}

	return nil
}

func (a Storable) decodeLog(log web3types.Log, event string) (*AuctionEvent, error) {
	var d AuctionEvent

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	var decodedData interface{}

	switch event {
	case AuctionCreated:
		var decoded LogAuctionCreated
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case ERC721Deposit:
		var decoded LogERC721Deposit
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case ERC721Withdrawal:
		var decoded LogERC721Withdrawal
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case BidSubmitted:
		var decoded LogBidSubmitted
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case BidWithdrawal:
		var decoded LogBidWithdrawal
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case BidMatched:
		var decoded LogBidMatched
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case AuctionExtended:
		var decoded LogAuctionExtended
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData=decoded
	case AuctionCanceled:
		var decoded LogAuctionCanceled
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case AuctionRevenueWithdrawal:
		var decoded LogAuctionRevenueWithdrawal
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case SlotRevenueCaptured:
		var decoded LogSlotRevenueCaptured
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	case ERC721RewardsClaim:
		var decoded LogERC721RewardsClaim
		err = a.auctionAbi.UnpackIntoInterface(&decoded, event, data)
		if err != nil {
			return nil, errors.Wrap(err, "could not unpack log data")
		}
		decodedData = decoded
	default:
		logger.Debug("Unknown auction event")
	}

	json, err := json.Marshal(decodedData)
	if err != nil {
		return nil, errors.Wrap(err, "could not pack data to json")
	}

	d.data = json

	d.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from kek contract to int64")
	}

	d.TransactionHash = log.TransactionHash
	d.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  kek contract to int64")
	}

	return &d, nil
}

func (a Storable) storeActions(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("auctions", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeErc721DepositEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("deposited_erc721", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeErc721WithdrawEvents (tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("withdrawn_erc721", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeAuctionCanceledEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("auctions_canceled", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeSubmittedBidEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("bids_submitted", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeWithdrawnBidEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("bids_withdrawn", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeAuctionExtendedEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("auctions_extended", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeBidMatchedEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("matched_bids", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeRevenueWithdrawEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("withdrawn_revenues", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeCaptureSlotRevenueEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("captured_revenues", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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

func (a Storable) storeErc721ClaimEvents(tx *sql.Tx, actions []AuctionEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("claimed_erc721", "tx_hash", "tx_index", "log_index", "data", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(a.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(a.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, a := range actions {
		_, err = stmt.Exec(a.TransactionHash, a.TransactionIndex, a.LogIndex, a.data, blockTimestamp, blockNumber)
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
