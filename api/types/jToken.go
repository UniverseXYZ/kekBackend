package types

type Transfer struct {
	LogData
	SYAddress string `json:"sy_address"`
	From      string `json:"sender"`
	To        string `json:"receiver"`
	Value     int64  `json:"value"`
}
