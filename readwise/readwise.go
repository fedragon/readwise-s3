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

type ListResponse struct {
	Count    uint64          `json:"count"`
	Next     *JsonURL        `json:"next"`
	Previous *JsonURL        `json:"previous"`
	Results  json.RawMessage `json:"results"`
}

type ListRequest struct {
	PageSize uint16 // default: 100, max: 1000
}

func Get(url *url.URL, token string, ghreq *ListRequest) (*ListResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	h := httpReq.Header
	h.Add("Authorization", "Token "+token)

	q := httpReq.URL.Query()
	q.Add("page_size", strconv.Itoa(int(ghreq.PageSize)))
	httpReq.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	httpRes, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	var res ListResponse
	if err = json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
