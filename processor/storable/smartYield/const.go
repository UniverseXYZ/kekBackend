package smartYield

const ZeroAddress = "0x0000000000000000000000000000000000000000"

const (
	BuyTokensEvent        = "BuyTokens"
	SellTokensEvent       = "SellTokens"
	BuySeniorBondEvent    = "BuySeniorBond"
	RedeemSeniorBondEvent = "RedeemSeniorBond"
	BuyJuniorBondEvent    = "BuyJuniorBond"
	RedeemJuniorBondEvent = "RedeemJuniorBond"
	TransferEvent         = "Transfer"
	HarvestEvent          = "Harvest"
	TransferFeesEvent     = "TransferFees"
)

type txType string

const (
	JuniorDeposit         txType = "JUNIOR_DEPOSIT"
	JuniorInstantWithdraw txType = "JUNIOR_INSTANT_WITHDRAW"
	JuniorRegularWithdraw txType = "JUNIOR_REGULAR_WITHDRAW"
	JuniorRedeem          txType = "JUNIOR_REDEEM"
	SeniorDeposit         txType = "SENIOR_DEPOSIT"
	SeniorRedeem          txType = "SENIOR_REDEEM"
	JtokenSend            txType = "JTOKEN_SEND"
	JtokenReceive         txType = "JTOKEN_RECEIVE"
	JbondSend             txType = "JBOND_SEND"
	JbondReceive          txType = "JBOND_RECEIVE"
	SbondSend             txType = "SBOND_SEND"
	SbondReceive          txType = "SBOND_RECEIVE"
)
