package be

type Transaction struct {
	Hash        string `json:"hash"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Amount      int64  `json:"amount"`
	Fee         int64  `json:"fee"`
	Parent      struct {
		Hash    string `json:"hash"`
		Ordinal int64  `json:"ordinal"`
	} `json:"parent"`
	Salt            int64  `json:"salt"`
	BlockHash       string `json:"blockHash"`
	SnapshotHash    string `json:"snapshotHash"`
	SnapshotOrdinal int64  `json:"snapshotOrdinal"`
}
