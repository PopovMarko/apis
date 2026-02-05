package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
		expOut string
		resp   struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{name: "Results",
			expErr: nil,
			expOut: "-   1   Task_1\n-   2   Task_2\n",
			resp:   testServerResponse["resultsMany"],
		},
		{name: "NoResults",
			expErr: ErrInvalid,
			resp:   testServerResponse["noResults"],
		},
		{name: "Invalid URL",
			expErr:      ErrConnection,
			resp:        testServerResponse["noResults"],
			closeServer: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanUp := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})
			defer cleanUp()
			if tc.closeServer {
				cleanUp()
			}
			var out bytes.Buffer
			err := listAction(&out, url)
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error: %s, got %s", tc.expErr, err)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %s, got %s", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expect NO error, got %s", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected out %q, got %q", tc.expOut, out.String())
			}
		})

	}
}

func TestViewAction(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
		expOut string
		resp   struct {
			Status int
			Body   string
		}
		id string
	}{
		{
			name:   "One Result",
			expErr: nil,
			expOut: "",
			resp:   testServerResponse["resultOne"],
			id:     "1",
		},
		{
			name:   "Not Found",
			expErr: ErrNotFound,
			expOut: "404 - not found",
			resp:   testServerResponse["notFound"],
			id:     "1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanUp := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})
			defer cleanUp()
			var out bytes.Buffer
			err := viewAction(&out, url, tc.id)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %s, got nil", tc.expErr)
					return
				}
				if errors.Is(tc.expErr, err) {
					t.Errorf("Expected error: %s, got %s", tc.expErr, err)
					return
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
		})
	}
}
