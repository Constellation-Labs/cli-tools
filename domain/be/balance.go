package be

type Balance struct {
	Address         string `json:"address"`
	Balance         int64  `json:"balance"`
	SnapshotHash    string `json:"snapshotHash"`
	SnapshotOrdinal int64  `json:"snapshotOrdinal"`
}
