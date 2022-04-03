package elasticsearch

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
)

type ElasticSearch struct {
	Client *elasticsearch.Client
	Index  string
	Alias  string
}

func New(addresses []string) (*ElasticSearch, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticSearch{
		Client: client,
	}, nil
}

func (e *ElasticSearch) CreateIndex(index, alias string) error {
	e.Index = index
	e.Alias = alias

	res, err := e.Client.Indices.Exists([]string{e.Index})
	if err != nil {
		return fmt.Errorf("cannot check index existence: %w", err)
	}

	if res.StatusCode == http.StatusOK {
		return nil
	}

	if res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error in index existence response: %s", res.String())
	}

	res, err = e.Client.Indices.Create(e.Index)
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error in index creation response: %s", res.String())
	}

	res, err = e.Client.Indices.PutAlias([]string{e.Index}, e.Alias)
	if err != nil {
		return fmt.Errorf("cannot create index alias: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error in index alias creation response: %s", res.String())
	}

	return nil
}
