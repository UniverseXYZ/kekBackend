package types

type ERC721Transfer struct {
	LogData

	TokenAddress string
	TokenType    string
	Sender       string
	Receiver     string
	TokenID      int64
}
