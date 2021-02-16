package smartYield

import (
	"math/big"

	"github.com/barnbridge/barnbridge-backend/types"
)

type TokenBuyTrade struct {
	*types.Event

	BuyerAddress string
	UnderlyingIn *big.Int
	TokensOut    *big.Int
	Fee          *big.Int
}

type TokenSellTrade struct {
	*types.Event

	SellerAddress string
	TokensIn      *big.Int
	UnderlyingOut *big.Int
	Forfeits      *big.Int
}

type SeniorBondBuyTrade struct {
	*types.Event

	BuyerAddress string
	SeniorBondID *big.Int
	UnderlyingIn *big.Int
	Gain         *big.Int
	ForDays      *big.Int
}

type SeniorBondRedeemTrade struct {
	*types.Event

	OwnerAddress string
	SeniorBondID *big.Int
	Fee          *big.Int
}

type JuniorBondBuyTrade struct {
	*types.Event

	BuyerAddress string
	JuniorBondID *big.Int
	TokensIn     *big.Int
	MaturesAt    *big.Int
}

type JuniorBondRedeemTrade struct {
	*types.Event

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

type SmartBondTransfer struct {
	*types.Event

	TokenAddress string
	Sender       string
	Receiver     string
	TokenID      *big.Int
}
