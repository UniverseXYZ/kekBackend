package types

type SYPool struct {
	ProtocolId         string `json:"protocolId"`
	ControllerAddress  string `json:"controllerAddress"`
	ModelAddress       string `json:"modelAddress"`
	ProviderAddress    string `json:"providerAddress"`
	SmartYieldAddress  string `json:"smartYieldAddress"`
	OracleAddress      string `json:"oracleAddress"`
	JuniorBondAddress  string `json:"juniorBondAddress"`
	SeniorBondAddress  string `json:"seniorBondAddress"`
	CTokenAddress      string `json:"cTokenAddress"`
	UnderlyingAddress  string `json:"underlyingAddress"`
	UnderlyingSymbol   string `json:"underlyingSymbol"`
	UnderlyingDecimals int64  `json:"underlyingDecimals"`
}
