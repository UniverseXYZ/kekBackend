package smartYield

import (
	"math/big"
)

type TokenBuyTrade struct {
	BuyerAddress string
	UnderlyingIn *big.Int
	TokensOut    *big.Int
	Fee          *big.Int

	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type TokenSellTrade struct {
	SellerAddress string
	TokensIn      *big.Int
	UnderlyingOut *big.Int
	Forfeits      *big.Int

	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type SeniorBondBuyTrade struct {
	BuyerAddress string
	SeniorBondID *big.Int
	UnderlyingIn *big.Int
	Gain         *big.Int
	ForDays      *big.Int

	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type SeniorBondRedeemTrade struct {
	OwnerAddress string
	SeniorBondID *big.Int
	Fee          *big.Int

	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type JuniorBondBuyTrade struct {
	BuyerAddress string
	JuniorBondID *big.Int
	TokensIn     *big.Int
	MaturesAt    *big.Int

	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type JuniorBondRedeemTrade struct {
	OwnerAddress  string
	JuniorBondID  *big.Int
	UnderlyingOut *big.Int

	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type TokenTrades struct {
	tokenBuyTrades  []TokenBuyTrade
	tokenSellTrades []TokenSellTrade
}

type SeniorTrades struct {
	seniorBondRedeems []SeniorBondRedeemTrade
	seniorBondBuys    []SeniorBondBuyTrade
}

type JuniorTrades struct {
	juniorBondRedeems []JuniorBondRedeemTrade
	juniorBondBuys    []JuniorBondBuyTrade
}
