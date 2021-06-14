package universe

import (
	"database/sql"
	"encoding/hex"
	"math/big"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kekDAO/kekBackend/state"
	"github.com/kekDAO/kekBackend/types"
	"github.com/kekDAO/kekBackend/utils"
)

var logger = logrus.WithField("module", "storable(universe)")

type Storable struct {
	config      Config
	raw         *types.RawData
	universeAbi abi.ABI
}

type MintedEvent struct {
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
	TokenID          string
	TokenURI         string
}

type DeployedEvent struct {
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
	TokenName        string
	TokenSymbol      string
	ContractAddress  string
}

type LogUniverseERC721ContractDeployed struct {
	TokenName       string         "json:\"tokenName\""
	TokenSymbol     string         "json:\"tokenSymbol\""
	ContractAddress common.Address "json:\"contractAddress\""
	Time            *big.Int       "json:\"time\""
}

type UniverseERC721TokenMinted struct {
	TokenID  *big.Int "json:\"tokenId\""
	TokenURI string   "json:\"tokenURI\""
	Time     *big.Int "json:\"time\""
}

func NewStorable(config Config, raw *types.RawData, universeAbi abi.ABI) *Storable {
	return &Storable{
		config:      config,
		raw:         raw,
		universeAbi: universeAbi,
	}
}

func (u *Storable) ToDB(tx *sql.Tx) error {
	var deployedEvents []DeployedEvent
	var mintedEvents []MintedEvent

	for _, data := range u.raw.Receipts {
		for _, log := range data.Logs {
			if utils.CleanUpHex(log.Address) != utils.CleanUpHex(u.config.Address) {
				continue
			}
			if len(log.Topics) == 0 {
				continue
			}

			if utils.LogIsEvent(log, u.universeAbi, Deployed) {
				d, err := u.decodeLog(log, Deployed)
				if err != nil {
					return err
				}

				deployedEvents = append(deployedEvents, *d)

				state.AddMonitoredNFTToState(d.ContractAddress)
				state.AddMonitoredNFTToDB(d.ContractAddress)
			}

			if utils.LogIsEvent(log, u.universeAbi, Minted) &&
				(state.IsMonitoredAccount(log) || u.IsPublicCollection(log)) {
				m, err := u.decodeMintedLog(log, Minted)
				if err != nil {
					return err
				}

				mintedEvents = append(mintedEvents, *m)

			}
		}
	}

	if len(deployedEvents) > 0 {
		err := u.storeEvents(tx, deployedEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "deployed events").Debug("no events found")
	}

	if len(mintedEvents) > 0 {
		err := u.storeMintedEvents(tx, mintedEvents)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("handler", "minted events").Debug("no events found")
	}

	return nil
}

func (u Storable) decodeLog(log web3types.Log, event string) (*DeployedEvent, error) {
	var d DeployedEvent

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	var decoded LogUniverseERC721ContractDeployed
	err = u.universeAbi.UnpackIntoInterface(&decoded, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	d.TokenName = decoded.TokenName
	d.TokenSymbol = decoded.TokenSymbol
	d.ContractAddress = decoded.ContractAddress.String()

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

func (u Storable) decodeMintedLog(log web3types.Log, event string) (*MintedEvent, error) {
	var m MintedEvent

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	var decoded UniverseERC721TokenMinted
	err = u.universeAbi.UnpackIntoInterface(&decoded, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	m.TokenID = decoded.TokenID.String()
	m.TokenURI = decoded.TokenURI

	m.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from kek contract to int64")
	}

	m.TransactionHash = log.TransactionHash
	m.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  kek contract to int64")
	}

	return &m, nil
}

func (u Storable) IsPublicCollection(log web3types.Log) bool {
	if len(log.Topics) >= 3 {
		if utils.NormalizeAddress(u.config.PublicCollection) == utils.Topic2Address(log.Topics[1]) ||
			utils.NormalizeAddress(u.config.PublicCollection) == utils.Topic2Address(log.Topics[2]) {
			return true
		}
	}
	return false
}

func (u Storable) storeEvents(tx *sql.Tx, events []DeployedEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("universe", "tx_hash", "tx_index", "log_index", "token_name", "token_symbol", "contract_address", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(u.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(u.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, e := range events {
		_, err = stmt.Exec(e.TransactionHash, e.TransactionIndex, e.LogIndex, e.TokenName, e.TokenSymbol, e.ContractAddress, blockTimestamp, blockNumber)
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

func (u Storable) storeMintedEvents(tx *sql.Tx, events []MintedEvent) error {
	stmt, err := tx.Prepare(pq.CopyIn("minted_nfts", "tx_hash", "tx_index", "log_index", "token_id", "token_uri", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseInt(u.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	blockTimestamp, err := strconv.ParseInt(u.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	for _, e := range events {
		_, err = stmt.Exec(e.TransactionHash, e.TransactionIndex, e.LogIndex, e.TokenID, e.TokenURI, blockTimestamp, blockNumber)
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
