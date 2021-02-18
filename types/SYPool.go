package types

type SYPool struct {
	ProtocolId          string
	ControllerAddress   string
	ModelAddress        string
	ProviderAddress     string
	SmartYieldAddress   string
	OracleAddress       string
	JuniorBondAddress   string
	SeniorBondAddress   string
	ReceiptTokenAddress string
	UnderlyingAddress   string
	UnderlyingSymbol    string
	UnderlyingDecimals  int64
}
