package smartYield

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

type CompoundProvider struct {
	transfersFees []TransferFees
}

type CompoundController struct {
	harvests []Harvest
}
