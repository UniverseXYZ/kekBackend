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

type TxType string

const (
	JuniorDeposit         TxType = "JUNIOR_DEPOSIT"
	JuniorInstantWithdraw TxType = "JUNIOR_INSTANT_WITHDRAW"
	JuniorRegularWithdraw TxType = "JUNIOR_REGULAR_WITHDRAW"
	JuniorRedeem          TxType = "JUNIOR_REDEEM"
	SeniorDeposit         TxType = "SENIOR_DEPOSIT"
	SeniorRedeem          TxType = "SENIOR_REDEEM"
	JtokenSend            TxType = "JTOKEN_SEND"
	JtokenReceive         TxType = "JTOKEN_RECEIVE"
	JbondSend             TxType = "JBOND_SEND"
	JbondReceive          TxType = "JBOND_RECEIVE"
	SbondSend             TxType = "SBOND_SEND"
	SbondReceive          TxType = "SBOND_RECEIVE"
)
