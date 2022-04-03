package advertsvc

import "net/url"

type CreateRequest struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type UpdateRequest struct {
	ID    string
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type DeleteRequest struct {
	ID string
}

type SearchOneRequest struct {
	ID string
}

type SearchRequest struct {
	Values url.Values
}
