package advertsvc

import "time"

type CreateResponse struct {
	ID string `json:"id"`
}

type SearchResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
