package readwise

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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

type GetHighlightsResponse struct {
	Count    uint64          `json:"count"`
	Next     *JsonURL        `json:"next"`
	Previous *JsonURL        `json:"previous"`
	Results  json.RawMessage `json:"results"`
}

type GetHighlightsRequest struct {
	PageSize uint16 // default: 100, max: 1000
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
