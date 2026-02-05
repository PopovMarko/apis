package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"pragprog.com/rggo/interacting/todo"
)

var res struct {
	Results      todo.List `json:"results"`
	Date         int64     `json:"date"`
	TotalResults int       `json:"total_results"`
}

func setUpAPI(t *testing.T, listInit bool) (string, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp(".", "temp_todo.json")
	if err != nil {
		t.Fatal(err)
	}

	list := todo.NewList()
	if listInit {
		list.Add("Test task 1")
		list.Add("Test task 2")
	}
	list.Save(tempFile.Name())

	testS := httptest.NewServer(newMux(tempFile.Name()))

	return testS.URL, func() {
		testS.Close()
		os.Remove(tempFile.Name())
	}

}

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}
func TestGet(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{name: "Get root",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "There is rootHandler response",
		},
		{name: "Get all",
			path:       "/todo",
			expCode:    http.StatusOK,
			expItems:   2,
			expContent: "Test task 1",
		},
		{name: "Get one",
			path:       "/todo/1",
			expCode:    http.StatusOK,
			expItems:   1,
			expContent: "Test task 1",
		},
		{name: "Not found path",
			path:    "/index/500",
			expCode: http.StatusNotFound,
		},
	}

	url, cleanup := setUpAPI(t, true)
	defer cleanup()
	var tss []*http.Response

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)

			resp, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}

			tss = append(tss, resp)

			if resp.StatusCode != tc.expCode {
				t.Error(err)
			}

			switch {
			case strings.Contains(resp.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(resp.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Expect: %s, got: %s", tc.expContent, string(body))
				}
			case resp.Header.Get("Content-Type") == "application/json":
				if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
					t.Error(err)
				}
				resp.Body.Close()
				if res.TotalResults != tc.expItems {
					t.Errorf("Expect %d items, got %d", tc.expItems, res.TotalResults)
				}
				if res.Results[0].Task != tc.expContent {
					t.Errorf("Expect %s, got %s", tc.expContent, res.Results[0].Task)
				}

			default:
				t.Fatalf("Unsupported Content-Type: %q", resp.Header.Get("Content-Type"))
			}
		},
		)
	}

	defer func() {
		for _, s := range tss {
			s.Body.Close()
		}
	}()
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		expCode     int
		expItems    int
		expContent  string
		expResponse string
	}{
		{name: "Post new item",
			path:        "/todo",
			expCode:     http.StatusCreated,
			expItems:    1,
			expContent:  "Test task 1",
			expResponse: "Item added",
		},
		{name: "Post new item 2",
			path:        "/todo",
			expCode:     http.StatusCreated,
			expItems:    2,
			expContent:  "Test task 1",
			expResponse: "Item added",
		},
	}
	var tss []*http.Response

	url, cleanUp := setUpAPI(t, false)
	defer cleanUp()

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := struct {
				Task string `json:"task"`
			}{
				Task: fmt.Sprintf("Test task %d", i+1),
			}

			var body bytes.Buffer
			json.NewEncoder(&body).Encode(&item)

			resp, err := http.Post(url+tc.path, "application/json", &body)
			if err != nil {
				t.Fatal(err)
			}
			tss = append(tss, resp)
			ans, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(ans) != tc.expResponse {
				t.Errorf("Expected reply: %s, got %s", tc.expResponse, string(ans))
			}

			if resp.StatusCode != tc.expCode {
				t.Errorf("Expected: %d, got %d", tc.expCode, resp.StatusCode)
			}
		},
		)
		t.Run("Check "+tc.name, func(t *testing.T) {
			respCk, err := http.Get(url + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			tss = append(tss, respCk)
			if err := json.NewDecoder(respCk.Body).Decode(&res); err != nil {
				t.Fatal(err)
			}
			if len(res.Results) != tc.expItems {
				t.Errorf("Should be %d items, got %d", tc.expItems, len(res.Results))
			}
			if res.Results[0].Task != tc.expContent {
				t.Errorf("Task name should be %s, got %s", tc.expContent, res.Results[0].Task)
			}
		})
	}
	defer func() {
		for _, r := range tss {
			r.Body.Close()
		}
	}()
}

func TestDelete(t *testing.T) {
	url, cleanUp := setUpAPI(t, true)
	defer cleanUp()

	t.Run("Delete", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, url+"/todo/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status: %d, got %d", http.StatusNoContent, resp.StatusCode)
		}
	})
	t.Run("Check delete", func(t *testing.T) {
		r, err := http.Get(url + "/todo")
		if err != nil {
			t.Fatal(err)
		}

		var resp todoResponse

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if len(resp.Results) != 1 {
			t.Errorf("Expect 1 item left, got %d", len(resp.Results))
		}
		if resp.Results[0].Task != "Test task 2" {
			t.Errorf("Expect 'Test task 2' left, got %q", resp.Results[0].Task)
		}
	})
}

func TestComplete(t *testing.T) {
	url, cleanUp := setUpAPI(t, true)
	defer cleanUp()

	t.Run("Complete", func(t *testing.T) {
		q, err := http.NewRequest(http.MethodPatch, url+"/todo/2?complete", nil)
		if err != nil {
			t.Fatal(err)
		}
		r, err := http.DefaultClient.Do(q)
		if err != nil {
			t.Fatal(err)
		}
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected status: %d, got %d", http.StatusOK, r.StatusCode)
		}
		resp, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		if string(resp) != "Item status changed" {
			t.Errorf("Expected response 'Item status changed', got %s", string(resp))
		}
	})
	t.Run("Check complete", func(t *testing.T) {
		r, err := http.Get(url + "/todo")
		if err != nil {
			t.Fatal(err)
		}

		resp := todoResponse{}
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if !resp.Results[1].Done {
			t.Errorf("Expect item Done is true, got %t", resp.Results[1].Done)
		}
	})
}
