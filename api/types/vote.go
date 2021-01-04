package types

type Vote struct {
	User           string `json:"address"`
	Support        bool   `json:"support"`
	BlockTimestamp int64  `json:"blockTimestamp"`
	Power          string `json:"power"`
}
