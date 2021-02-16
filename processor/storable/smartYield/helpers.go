package smartYield

import (
	"encoding/hex"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (s *Storable) decodeTokenBuyEvent(log web3types.Log, event string) (*TokenBuyTrade, error) {
	var t TokenBuyTrade

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.BuyerAddress = utils.Topic2Address(log.Topics[1])

	return &t, nil
}

func (s *Storable) decodeTokenSellEvent(log web3types.Log, event string) (*TokenSellTrade, error) {
	var t TokenSellTrade

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.SellerAddress = utils.Topic2Address(log.Topics[1])

	return &t, nil
}

func (s *Storable) decodeSeniorBondBuyEvent(log web3types.Log, event string) (*SeniorBondBuyTrade, error) {
	var t SeniorBondBuyTrade

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.BuyerAddress = utils.Topic2Address(log.Topics[1])

	n := new(big.Int)
	n, ok := n.SetString(log.Topics[2], 10)
	if !ok {
		return nil, errors.New("could not convert seniorBondId ")
	}

	t.SeniorBondID = n

	return &t, nil
}

func (s *Storable) decodeSeniorBondRedeemEvent(log web3types.Log, event string) (*SeniorBondRedeemTrade, error) {
	var t SeniorBondRedeemTrade

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.OwnerAddress = utils.Topic2Address(log.Topics[1])

	n := new(big.Int)
	n, ok := n.SetString(log.Topics[2], 10)
	if !ok {
		return nil, errors.New("could not convert seniorBondId ")
	}

	t.SeniorBondID = n

	return &t, nil
}

func (s *Storable) decodeJuniorBondBuyEvent(log web3types.Log, event string) (*JuniorBondBuyTrade, error) {
	var t JuniorBondBuyTrade

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.BuyerAddress = utils.Topic2Address(log.Topics[1])

	n := new(big.Int)
	n, ok := n.SetString(log.Topics[2], 10)
	if !ok {
		return nil, errors.New("could not convert JuniorBondID ")
	}

	t.JuniorBondID = n

	return &t, nil
}

func (s *Storable) decodeJuniorBondRedeemEvent(log web3types.Log, event string) (*JuniorBondRedeemTrade, error) {
	var t JuniorBondRedeemTrade

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.OwnerAddress = utils.Topic2Address(log.Topics[1])

	n := new(big.Int)
	n, ok := n.SetString(log.Topics[2], 10)
	if !ok {
		return nil, errors.New("could not convert JuniorBondID ")
	}

	t.JuniorBondID = n

	return &t, nil
}
