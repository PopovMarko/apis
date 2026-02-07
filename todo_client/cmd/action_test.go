//go:build !integration

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

func TestAddAction(t *testing.T) {
	testCases := []struct {
		name           string
		expUrlPath     string
		expMethod      string
		expBody        string
		expContentType string
		expErr         error
		expOut         string
		args           []string
		resp           struct {
			Status int
			Body   string
		}
	}{
		{name: "Add request",
			expUrlPath:     "/todo",
			expMethod:      "POST",
			expBody:        `{"Task":"Task 1"}` + "\n",
			expContentType: "application/json",
			expErr:         nil,
			expOut:         "Task: Task_1 added to the list",
			args:           []string{"Task", "1"},
			resp:           testServerResponse["created"],
		},
		{name: "Add bad request",
			expUrlPath:     "/todo",
			expMethod:      "POST",
			expBody:        `{"Task":""}` + "\n",
			expContentType: "application/json",
			expErr:         ErrInvalidResponse,
			expOut:         "Task: Task_1 added to the list",
			args:           []string{""},
			resp:           testServerResponse["badRequest"],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanUp := mockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.expUrlPath {
					t.Errorf("Expected path: %s, got %s", tc.expUrlPath, r.URL.Path)
				}
				if r.Method != tc.expMethod {
					t.Errorf("Expected method: %s, got %s", tc.expMethod, r.Method)
				}
				body, err := io.ReadAll(r.Body)

				if err != nil {
					t.Fatalf("Can't read the Body: %s", err)
				}
				if string(body) != tc.expBody {
					t.Errorf("Expected body: %s, got %s", tc.expBody, string(body))
				}
				contentType := r.Header.Get("Content-Type")
				if contentType != tc.expContentType {
					t.Errorf("Expected content-type: %s, got %s", tc.expContentType, contentType)
				}
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})

			defer cleanUp()
			var out bytes.Buffer
			err := addAction(&out, url, tc.args)
			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %s, got %s", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpeted error: %s", err)
			}
		},
		)
	}
}

func TestCompleteAction(t *testing.T) {
	testCases := []struct {
		name       string
		expUrlPath string
		expMethod  string
		expErr     error
		expOut     string
		args       []string
		resp       struct {
			Status int
			Body   string
		}
	}{
		{name: "Complete",
			expUrlPath: "/todo/1",
			expMethod:  "PATCH",
			expErr:     nil,
			expOut:     "Item No 1 set as completed",
			args:       []string{"1"},
			resp:       testServerResponse["noContent"],
		},
		{name: "Complete without arg",
			expUrlPath: "/todo/",
			expMethod:  "PATCH",
			expErr:     ErrNotNumber,
			expOut:     "",
			args:       []string{""},
			resp:       testServerResponse["badRequest"],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanUp := mockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.expUrlPath {
					t.Errorf("Expected path: %s, got %s", tc.expUrlPath, r.URL.Path)
				}
				if r.Method != tc.expMethod {
					t.Errorf("Expected method: %s, got %s", tc.expMethod, r.Method)
				}
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)

			},
			)
			defer cleanUp()
			var out bytes.Buffer
			err := completeAction(&out, url, tc.args)
			if tc.expErr != nil {
				if !errors.Is(tc.expErr, errors.Unwrap(err)) {
					t.Errorf("Expected error %s, got %s", tc.expErr, errors.Unwrap(err))
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}
		},
		)
	}
}

func TestRemoveAction(t *testing.T) {
	testCases := []struct {
		name       string
		expUrlPath string
		expMethod  string
		expErr     error
		expOut     string
		args       []string
		resp       struct {
			Status int
			Body   string
		}
	}{
		{name: "Delete",
			expUrlPath: "/todo/1",
			expMethod:  "DELETE",
			expErr:     nil,
			expOut:     "Item No 1 deleted",
			args:       []string{"1"},
			resp:       testServerResponse["noContent"],
		},
		{name: "Delete without arg",
			expUrlPath: "/todo/",
			expMethod:  "DELETE",
			expErr:     ErrNotNumber,
			expOut:     "",
			args:       []string{""},
			resp:       testServerResponse["badRequest"],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanUp := mockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.expUrlPath {
					t.Errorf("Expected path: %s, got %s", tc.expUrlPath, r.URL.Path)
				}
				if r.Method != tc.expMethod {
					t.Errorf("Expected method: %s, got %s", tc.expMethod, r.Method)
				}
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)

			},
			)
			defer cleanUp()
			var out bytes.Buffer
			err := removeAction(&out, url, tc.args)
			if tc.expErr != nil {
				if !errors.Is(tc.expErr, errors.Unwrap(err)) {
					t.Errorf("Expected error %s, got %s", tc.expErr, errors.Unwrap(err))
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}
		},
		)
	}
}
