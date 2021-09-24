package main

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_CanUnmarshalJSONResponse(t *testing.T) {
	input := `
	{
		"count": 1163,
		"next": "https://readwise.io/api/v2/highlights?page=2",
		"previous": null,
		"results": [
			{
				"id": 59758950,
				"text": "The fundamental belief of metaphysicians is THE BELIEF IN ANTITHESES OF VALUES.",
				"note": "",
				"location": 9,
				"location_type": "order",
				"highlighted_at": null,
				"url": null,
				"color": "",
				"updated": "2020-10-01T12:58:44.716235Z",
				"book_id": 2608248,
				"tags": [
					{
						"id": 123456,
						"name": "philosophy"
					}
				]
			}
		]
	}`

	response := GetHighlightsResponse{}
	if err := json.Unmarshal([]byte(input), &response); err != nil {
		t.Errorf("expected success, but got error: %v", err)
	}

	next, _ := url.Parse("https://readwise.io/api/v2/highlights?page=2")
	updated, _ := time.Parse(time.RFC3339Nano, "2020-10-01T12:58:44.716235Z")
	expected := GetHighlightsResponse{
		Count:    1163,
		Next:     &JsonURL{Url: *next},
		Previous: nil,
		Results: []Highlight{{
			ID:            59758950,
			Text:          "The fundamental belief of metaphysicians is THE BELIEF IN ANTITHESES OF VALUES.",
			Note:          "",
			Location:      9,
			LocationType:  "order",
			HighlightedAt: nil,
			URL:           nil,
			Color:         "",
			Updated:       updated,
			BookID:        2608248,
			Tags: []Tag{{
				ID:   123456,
				Name: "philosophy",
			}},
		}},
	}

	if !reflect.DeepEqual(response, expected) {
		t.Errorf("expected %v, but got %v", expected, response)
	}
}
