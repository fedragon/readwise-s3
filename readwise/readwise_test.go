package readwise

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecodesResponseOn200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{ "count": 0, "next": "http://localhost?page=3", "previous": "http://localhost?page=1", "results": []}`))
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	ctx := NewListRequestContext(server.Client(), u, "xyz", 10)
	res, err := Get(ctx)
	if err != nil {
		t.Errorf("expected result but got %v", err)
	}

	assert.Equal(t, uint64(0), res.Count)
	assert.Equal(t, "http://localhost?page=3", res.Next.Url.String())
	assert.Equal(t, "http://localhost?page=1", res.Previous.Url.String())
}

func Test_RetriesWithDelayOnTooManyRequests(t *testing.T) {
	respondWith429 := true
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if respondWith429 {
			respondWith429 = false
			rw.Header().Add("Retry-After", "1")
			rw.WriteHeader(http.StatusTooManyRequests)
		} else {
			rw.Write([]byte(`{ "count": 0, "next": null, "previous": null, "results": []}`))
			rw.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	ctx := NewListRequestContext(server.Client(), u, "xyz", 10)
	_, err := Get(ctx)
	assert.Nil(t, err, "Unexpected error")
}

func Test_EventuallyStopsRetrying(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Retry-After", "1")
		rw.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	ctx := NewListRequestContext(server.Client(), u, "xyz", 10)
	_, err := Get(ctx)
	assert.NotNil(t, err, "Unexpected error")
}
