package smartYield

const (
	BUY_TOKENS_EVENT         = "BuyTokens"
	SELL_TOKENS_EVENT        = "SellTokens"
	BUY_SENIOR_BOND_EVENT    = "BuySeniorBond"
	REDEEM_SENIOR_BOND_EVENT = "RedeemSeniorBond"
	BUY_JUNIOR_BOND_EVENT    = "BuyJuniorBond"
	REDEEM_JUNIOR_BOND_EVENT = "RedeemJuniorBond"
	TRANSFER_EVENT           = "Transfer"
	HARVEST_EVENT            = "Harvest"
	TRANSFER_FEES_EVENT      = "TransferFees"
)

type txType string

const (
	JUNIOR_DEPOSIT          txType = "JUNIOR_DEPOSIT"
	JUNIOR_INSTANT_WITHDRAW txType = "JUNIOR_INSTANT_WITHDRAW"
	JUNIOR_REGULAR_WITHDRAW txType = "JUNIOR_REGULAR_WITHDRAW"
	JUNIOR_REDEEM           txType = "JUNIOR_REDEEM"
	SENIOR_DEPOSIT          txType = "SENIOR_DEPOSIT"
	SENIOR_REDEEM           txType = "SENIOR_REDEEM"
	JTOKEN_TRANSFER         txType = "JTOKEN_TRANSFER"
	JBOND_TRANSFER          txType = "JBOND_TRANSFER"
	SBOND_TRANSFER          txType = "SBOND_TRANSFER"
)
