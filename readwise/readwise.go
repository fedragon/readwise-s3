package readwise

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

type ListResponse struct {
	Count    uint64          `json:"count"`
	Next     *JsonURL        `json:"next"`
	Previous *JsonURL        `json:"previous"`
	Results  json.RawMessage `json:"results"`
}

type ListRequestContext struct {
	client       *http.Client
	url          *url.URL
	token        string
	pageSize     uint16 // default: 100, max: 1000
	attemptsLeft uint8
}

func NewListRequestContext(client *http.Client, url *url.URL, token string, pageSize uint16) *ListRequestContext {
	return &ListRequestContext{
		client:       client,
		url:          url,
		token:        token,
		pageSize:     pageSize,
		attemptsLeft: 3,
	}
}

func (r *ListRequestContext) SetURL(newURL *url.URL) {
	r.url = newURL
}

func Get(ctx *ListRequestContext) (*ListResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, ctx.url.String(), nil)
	if err != nil {
		return nil, err
	}

	h := httpReq.Header
	h.Add("Authorization", "Token "+ctx.token)

	q := httpReq.URL.Query()
	q.Add("page_size", strconv.Itoa(int(ctx.pageSize)))
	httpReq.URL.RawQuery = q.Encode()

	httpRes, err := ctx.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case http.StatusOK:
		var res ListResponse
		if err = json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}
		return &res, nil
	case http.StatusTooManyRequests:
		if ctx.attemptsLeft == 0 {
			return nil, errors.New("reached maximum number of retries, giving up")
		}
		retryAfter := httpRes.Header.Get("Retry-After")
		if retryAfter != "" {
			seconds, err := strconv.Atoi(retryAfter)
			if err != nil {
				return nil, fmt.Errorf("unexpected value in Retry-After header: %v, error: %v", retryAfter, err)
			}
			fmt.Printf("Got %s. Waiting for %vs before retrying\n", httpRes.Status, seconds)
			ctx.attemptsLeft--

			<-time.After(time.Second * time.Duration(seconds))

			return Get(ctx)
		}

		return nil, errors.New("unexpected value in Retry-After header: empty string")
	default:
		return nil, fmt.Errorf("got unexpected HTTP response status: %v", httpRes.Status)
	}
}
