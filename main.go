package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

type GetHighlightsRequest struct {
	PageSize        uint16 // default: 100, max: 1000
	Page            uint64
	BookID          *uint64
	UpdatedLT       *time.Time
	UpdatedGT       *time.Time
	HighlightedAtLT *time.Time
	HighlightedAtGT *time.Time
}

func Get(url *url.URL, token string, ghreq *GetHighlightsRequest) (*GetHighlightsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	h := req.Header
	h.Add("Authorization", "Token "+token)

	q := req.URL.Query()
	q.Add("page_size", strconv.Itoa(int(ghreq.PageSize)))
	q.Add("page", strconv.Itoa(int(ghreq.Page)))
	if ghreq.BookID != nil {
		q.Add("book_id", strconv.Itoa(int(*ghreq.BookID)))
	}
	if ghreq.UpdatedLT != nil {
		q.Add("updated__lt", ghreq.UpdatedLT.Format(time.RFC3339Nano))
	}
	if ghreq.UpdatedGT != nil {
		q.Add("updated__gt", ghreq.UpdatedGT.Format(time.RFC3339Nano))
	}
	if ghreq.HighlightedAtLT != nil {
		q.Add("highlighted_at__lt", ghreq.HighlightedAtLT.Format(time.RFC3339Nano))
	}
	if ghreq.HighlightedAtGT != nil {
		q.Add("highlighted_at__gt", ghreq.HighlightedAtGT.Format(time.RFC3339Nano))
	}
	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ghres GetHighlightsResponse
	if err = json.NewDecoder(res.Body).Decode(&ghres); err != nil {
		return nil, err
	}

	return &ghres, nil
}

func main() {
	req := GetHighlightsRequest{
		PageSize: 10,
		Page:     1,
	}

	u, _ := url.Parse("https://readwise.io/api/v2/highlights/")

	res, err := Get(u, os.Getenv("API_TOKEN"), &req)
	if err != nil {
		log.Fatal(err)
	}

	for _, h := range res.Results {
		fmt.Println(h)
	}
}
