package be

type Snapshot struct {
	Hash             string   `json:"hash"`
	LastSnapshotHash string   `json:"lastSnapshotHash"`
	Ordinal          int64    `json:"ordinal"`
	Height           int64    `json:"height"`
	SubHeight        int64    `json:"subHeight"`
	Blocks           []string `json:"blocks"`
	Rewards          []struct {
		Amount      int64  `json:"amount"`
		Destination string `json:"destination"`
	} `json:"rewards"`
}
