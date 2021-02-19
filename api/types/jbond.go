package types

type JuniorBondBuyTrade struct {
	LogData

	SYAddress    string `json:"poolAddress"`
	BuyerAddress string `json:"buyerAddress"`
	JuniorBondID int64  `json:"juniorBondID"`
	TokensIn     int64  `json:"tokensIn"`
	MaturesAt    int64  `json:"maturesAt"`
}

type JuniorBondRedeemTrade struct {
	LogData

	SYAddress     string `json:"poolAddress"`
	OwnerAddress  string `json:"ownerAddress"`
	JuniorBondID  int64  `json:"juniorBondID"`
	UnderlyingOut int64  `json:"underlyingOut"`
}

type JuniorBondTrades struct {
	Buys      []JuniorBondBuyTrade
	Redeems   []JuniorBondRedeemTrade
	Transfers []ERC721Transfer
}
