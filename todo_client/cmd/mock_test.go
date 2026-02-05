package cmd

import (
	"net/http"
	"net/http/httptest"
)

var testServerResponse = map[string]struct {
	Status int
	Body   string
}{
	"resultsMany": {
		Status: http.StatusOK,
		Body: `{
			"results": [
			{
			"Task": "Task_1",
			"Done": false,
			"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
			"CompletedAt": "0001-01-01T00:00:00Z"
			},
			{
			"Task": "Task_2",
			"Done": false,
			"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
			"CompletedAt": "0001-01-01T00:00:00Z"
			}
			],
			"date": 1572265440,
			"total_results": 2
			}`,
	},

	"resultOne": {
		Status: http.StatusOK,
		Body: `{
			"results": [
			{
			"Task": "Task_1",
			"Done": false,
			"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
			"CompletedAt": "0001-01-01T00:00:00Z"
			}
			],
			"date": 1572265440,
			"total_results": 1
			}`,
	},

	"noResults": {
		Status: http.StatusOK,
		Body: `{
			"results": [],
			"date": 1572265440,
			"total_results": 0
			}`,
	},

	"root": {
		Status: http.StatusOK,
		Body:   "This is API handler respond",
	},

	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},
}

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() {
		ts.Close()
	}
}
