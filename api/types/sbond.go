package types

type SeniorBondBuyTrade struct {
	LogData

	SYAddress    string `json:"syAddress"`
	BuyerAddress string `json:"buyerAddress"`
	SeniorBondID int64  `json:"seniorBondId"`
	UnderlyingIn int64  `json:"underlyingIn"`
	Gain         int64  `json:"gain"`
	ForDays      int64  `json:"forDays"`
}

type SeniorBondRedeemTrade struct {
	LogData

	SYAddress    string `json:"syAddress"`
	OwnerAddress string `json:"ownerAddress"`
	SeniorBondID int64  `json:"seniorBondId"`
	Fee          int64  `json:"fee"`
}

type SeniorBondTrades struct {
	Buys      []SeniorBondBuyTrade
	Redeems   []SeniorBondRedeemTrade
	Transfers []ERC721Transfer
}
