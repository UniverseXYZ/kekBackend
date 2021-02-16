package smartYield

import (
	"math/big"
)

type TokenBuyTrade struct {
	BuyerAddress string
	UnderlyingIn *big.Int
	TokensOut    *big.Int
	Fee          *big.Int
}

type TokenSellTrade struct {
	SellerAddress string
	TokensIn      *big.Int
	UnderlyingOut *big.Int
	Forfeits      *big.Int
}

type SeniorBondBuyTrade struct {
	BuyerAddress string
	SeniorBondID *big.Int
	UnderlyingIn *big.Int
	Gain         *big.Int
	ForDays      *big.Int
}

type SeniorBondRedeemTrade struct {
	OwnerAddress string
	SeniorBondID *big.Int
	Fee          *big.Int
}

type JuniorBondBuyTrade struct {
	BuyerAddress string
	JuniorBondID *big.Int
	TokensIn     *big.Int
	MaturesAt    *big.Int
}

type JuniorBondRedeemTrade struct {
	OwnerAddress  string
	JuniorBondID  *big.Int
	UnderlyingOut *big.Int
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
