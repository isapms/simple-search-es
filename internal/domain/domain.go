package domain

import "time"

// SearchResponse represents generic response of Elasticsearch
type SearchResponse struct {
	Hits HitsResponse `json:"hits"`
}

// HitsResponse represents generic response of Elasticsearch Hits
type HitsResponse struct {
	Total    interface{} `json:"total"`
	MaxScore float64     `json:"max_score"`
	Hits     []Hit       `json:"hits"`
}

// Hit represents generic response of Elasticsearch Hit
type Hit struct {
	Index  string      `json:"_index"`
	Type   string      `json:"_type"`
	Id     string      `json:"_id"`
	Score  float64     `json:"_score"`
	Source interface{} `json:"_source"`
}

// Document represents a single document in Get API response body.
type Document struct {
	Source interface{} `json:"_source"`
}

// BoolQuery represents ES bool query structure
type BoolQuery struct {
	Filter []Term `json:"filter"`
}

// Term represents ES term query structure
type Term struct {
	Term map[string]interface{} `json:"term"`
}

// Advert represents advert model
type Advert struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Text      string     `json:"text"`
	Tags      []string   `json:"tags"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
