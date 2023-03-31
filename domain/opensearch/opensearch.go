package opensearch

import "dag-balances/domain/be"

type Response[T any] struct {
	Shards struct {
		Failed     int `json:"failed"`
		Skipped    int `json:"skipped"`
		Successful int `json:"successful"`
		Total      int `json:"total"`
	} `json:"_shards"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Hits []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type OpenSearch interface {
	GetTransactions(ordinal int64) ([]be.Transaction, error)
	GetBalance(address string, ordinal int64) (*be.Balance, error)
	PutBalance(balance be.Balance) error
}
