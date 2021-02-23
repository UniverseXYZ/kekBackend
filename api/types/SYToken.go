package types

type TokenBuyTrade struct {
	LogData

	SYAddress    string `json:"SYAddress"`
	BuyerAddress string `json:"buyerAddress"`
	UnderlyingIn int64  `json:"underlyingIn"`
	TokensOut    int64  `json:"tokensOut"`
	Fee          int64  `json:"fee"`
}

type TokenSellTrade struct {
	LogData

	SYAddress     string `json:"SYAddress"`
	SellerAddress string `json:"sellerAddress"`
	TokensIn      int64  `json:"tokensIn"`
	UnderlyingOut int64  `json:"underlyingOut"`
	Forfeits      int64  `json:"forfeits"`
}

type SYTokenTrades struct {
	Buys  []TokenBuyTrade
	Sells []TokenSellTrade
}
