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
	harvests      []Harvest
	transfersFees []TransferFees
}
