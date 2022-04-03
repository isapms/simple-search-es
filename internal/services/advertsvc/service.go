package advertsvc

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"simple-search-es/internal/domain"
	"simple-search-es/internal/repository/advertrepo"
	"time"
)

type Service struct {
	advertRepo advertrepo.Repository
}

func New(advertRepo advertrepo.Repository) Service {
	return Service{
		advertRepo: advertRepo,
	}
}

func (s Service) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	id := uuid.New().String()
	cr := time.Now().UTC()

	doc := domain.Advert{
		ID:        id,
		Title:     req.Title,
		Text:      req.Text,
		Tags:      req.Tags,
		CreatedAt: &cr,
	}

	if err := s.advertRepo.Insert(ctx, doc); err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{ID: id}, nil
}

func (s Service) Update(ctx context.Context, req UpdateRequest) error {
	doc := domain.Advert{
		ID:    req.ID,
		Title: req.Title,
		Text:  req.Text,
		Tags:  req.Tags,
	}

	if err := s.advertRepo.Update(ctx, doc); err != nil {
		return err
	}

	return nil
}

func (s Service) Delete(ctx context.Context, req DeleteRequest) error {
	if err := s.advertRepo.Delete(ctx, req.ID); err != nil {
		return err
	}

	return nil
}

func (s Service) SearchOne(ctx context.Context, req SearchOneRequest) (SearchResponse, error) {
	advert, err := s.advertRepo.FindOne(ctx, req.ID)
	if err != nil {
		return SearchResponse{}, err
	}

	return SearchResponse{
		ID:        advert.ID,
		Title:     advert.Title,
		Text:      advert.Text,
		Tags:      advert.Tags,
		CreatedAt: *advert.CreatedAt,
	}, nil
}

func (s Service) Search(ctx context.Context, request SearchRequest) (res []SearchResponse, err error) {
	boolQuery := domain.BoolQuery{
		Filter: []domain.Term{},
	}

	for key, value := range request.Values {
		term := domain.Term{
			Term: map[string]interface{}{key: value[0]},
		}
		boolQuery.Filter = append(boolQuery.Filter, term)
	}

	var myMap map[string]interface{}
	data, _ := json.Marshal(boolQuery)
	json.Unmarshal(data, &myMap)

	adverts, err := s.advertRepo.FindAll(ctx, myMap)
	if err != nil {
		return res, err
	}

	for _, advert := range adverts {
		res = append(res, SearchResponse{
			ID:        advert.ID,
			Title:     advert.Title,
			Text:      advert.Text,
			Tags:      advert.Tags,
			CreatedAt: *advert.CreatedAt,
		})
	}

	return res, nil
}
