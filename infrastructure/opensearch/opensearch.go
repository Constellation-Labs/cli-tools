package opensearch

import (
	"context"
	"dag-balances/domain/be"
	"dag-balances/domain/opensearch"
	"encoding/json"
	"fmt"
	"strings"

	os "github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type o struct {
	client            *os.Client
	transactionsIndex string
	balancesIndex     string
}

func (o *o) GetTransactions(ordinal int64) ([]be.Transaction, error) {

	content := strings.NewReader(fmt.Sprintf(`{
		"size": 1,
		"query": {
			"term": {
				"snapshotOrdinal": %d
			}
		}
	}`, ordinal))

	search := opensearchapi.SearchRequest{
		Index: []string{o.transactionsIndex},
		Body:  content,
	}

	searchResponse, err := search.Do(context.Background(), o.client)
	if err != nil {
		return nil, err
	}

	defer searchResponse.Body.Close()

	if searchResponse.IsError() {
		return nil, fmt.Errorf("An error occured when requesting transactions")
	}

	var response opensearch.Response[be.Transaction]
	if err := json.NewDecoder(searchResponse.Body).Decode(&response); err != nil {
		return nil, err
	}

	var txs []be.Transaction
	for _, el := range response.Hits.Hits {
		txs = append(txs, el.Source)
	}

	return txs, nil
}

func (o *o) GetBalance(address string, ordinal int64) (*be.Balance, error) {
	content := strings.NewReader(fmt.Sprintf(`{
		"size": 1,
		"query": {
			"bool": {
				"must": [
					{ "match": { "address": "%s" } }
				],
				"filter": [
					{ "term": { "snapshotOrdinal": %d } }
				]
			}
		}
	}`, address, ordinal))

	search := opensearchapi.SearchRequest{
		Index: []string{o.balancesIndex},
		Body:  content,
	}

	searchResponse, err := search.Do(context.Background(), o.client)
	if err != nil {
		return nil, err
	}

	defer searchResponse.Body.Close()

	if searchResponse.IsError() {
		return nil, fmt.Errorf("An error occured when requesting balance")
	}

	var response opensearch.Response[be.Balance]
	if err := json.NewDecoder(searchResponse.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.TimedOut {
		return nil, fmt.Errorf("Timed out")
	}

	if response.Shards.Successful == 0 {
		return nil, fmt.Errorf("No shards error")
	}

	if len(response.Hits.Hits) == 0 {
		return nil, nil
	}

	return &response.Hits.Hits[0].Source, nil
}

func (o *o) PutBalance(balance be.Balance) error {
	document := strings.NewReader(fmt.Sprintf(`{
		"address": "%s",
		"snapshotHash": "%s",
		"snapshotOrdinal": %d,
		"balance": %d
	}`, balance.Address, balance.SnapshotHash, balance.SnapshotOrdinal, balance.Balance))

	req := opensearchapi.IndexRequest{
		Index:      o.balancesIndex,
		DocumentID: fmt.Sprintf("%s%d", balance.Address, balance.SnapshotOrdinal),
		Body:       document,
	}

	_, err := req.Do(context.Background(), o.client)

	return err
}

func GetClient(url string) opensearch.OpenSearch {
	client, _ := os.NewClient(os.Config{
		Addresses: []string{url},
	})

	return &o{client: client, transactionsIndex: "transactions", balancesIndex: "balances"}
}

func GetDefaultClient() opensearch.OpenSearch {
	return GetClient("127.0.0.1")
}
