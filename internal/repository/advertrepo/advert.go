package advertrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"simple-search-es/elasticsearch"
	"simple-search-es/internal/domain"
	"simple-search-es/internal/utils"
	"time"
)

type Repository struct {
	elastic elasticsearch.ElasticSearch
	timeout time.Duration
}

func New(elastic elasticsearch.ElasticSearch) Repository {
	return Repository{
		elastic: elastic,
		timeout: time.Second * 10,
	}
}

func (a Repository) Insert(ctx context.Context, advert domain.Advert) error {
	body, err := json.Marshal(advert)
	if err != nil {
		return fmt.Errorf("insert: marshall: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      a.elastic.Alias,
		DocumentID: advert.ID,
		Body:       bytes.NewReader(body),
	}

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	res, err := req.Do(ctx, a.elastic.Client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return fmt.Errorf("error conflict")
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (a Repository) Update(ctx context.Context, advert domain.Advert) error {
	bdy, err := json.Marshal(advert)
	if err != nil {
		return fmt.Errorf("update: marshall: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      a.elastic.Alias,
		DocumentID: advert.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
	}

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	res, err := req.Do(ctx, a.elastic.Client)
	if err != nil {
		return fmt.Errorf("update: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return fmt.Errorf("error conflict")
	}

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}

func (a Repository) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      a.elastic.Alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	res, err := req.Do(ctx, a.elastic.Client)
	if err != nil {
		return fmt.Errorf("delete: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return fmt.Errorf("error conflict")
	}

	if res.IsError() {
		return fmt.Errorf("delete: response: %s", res.String())
	}

	return nil
}

func (a Repository) FindOne(ctx context.Context, id string) (domain.Advert, error) {
	req := esapi.GetRequest{
		Index:      a.elastic.Alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	res, err := req.Do(ctx, a.elastic.Client)
	if err != nil {
		return domain.Advert{}, fmt.Errorf("find one: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return domain.Advert{}, fmt.Errorf("error not found")
	}

	if res.IsError() {
		return domain.Advert{}, fmt.Errorf("find one: response: %s", res.String())
	}

	var advert domain.Advert
	var body domain.Document

	body.Source = &advert
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return domain.Advert{}, fmt.Errorf("find one: decode: %w", err)
	}

	return advert, nil
}

func (a Repository) FindAll(ctx context.Context, termQuery map[string]interface{}) ([]domain.Advert, error) {
	queryTerms := map[string]interface{}{
		"match_all": map[string]interface{}{},
	}
	if len(termQuery) > 0 {
		queryTerms = map[string]interface{}{
			"bool": termQuery,
		}
	}

	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": queryTerms,
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := a.elastic.Client.Search(
		a.elastic.Client.Search.WithIndex(a.elastic.Index),
		a.elastic.Client.Search.WithBody(&buf),
		a.elastic.Client.Search.WithTimeout(a.timeout),
		a.elastic.Client.Search.WithContext(ctx),
		a.elastic.Client.Search.WithTrackTotalHits(true),
		a.elastic.Client.Search.WithPretty(),
	)
	if err != nil {
		return []domain.Advert{}, fmt.Errorf("find one: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return []domain.Advert{}, fmt.Errorf("error not found")
	}

	if res.IsError() {
		return []domain.Advert{}, fmt.Errorf("find one: response: %s", res.String())
	}

	var r domain.SearchResponse

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return []domain.Advert{}, fmt.Errorf("find one: decode: %w", err)
	}

	var adverts []domain.Advert
	for _, hit := range r.Hits.Hits {
		advert := domain.Advert{}
		utils.MapToStruct(hit.Source, &advert)
		adverts = append(adverts, advert)
	}

	return adverts, nil
}
