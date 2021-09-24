package main

import (
	"encoding/json"
	"net/url"
	"time"
)

type JsonURL struct {
	Url url.URL
}

func (ju *JsonURL) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	ju.Url = *u

	return nil
}

type Tag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Highlight struct {
	ID            uint64     `json:"id"`
	Text          string     `json:"text"`
	Note          string     `json:"note"`
	Location      uint64     `json:"location"`
	LocationType  string     `json:"location_type"`
	HighlightedAt *time.Time `json:"highlighted_at"`
	URL           *JsonURL   `json:"url"`
	Color         string     `json:"color"`
	Updated       time.Time  `json:"updated"`
	BookID        uint64     `json:"book_id"`
	Tags          []Tag      `json:"tags"`
}

type GetHighlightsResponse struct {
	Count    uint64      `json:"count"`
	Next     *JsonURL    `json:"next"`
	Previous *JsonURL    `json:"previous"`
	Results  []Highlight `json:"results"`
}

func main() {

}
