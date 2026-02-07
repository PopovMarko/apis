package cmd

import (
	"bytes"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const timeFormat = "Jan/02 @15:04"

var (
	ErrConnection      = errors.New("Connection error")
	ErrNotFound        = errors.New("Not found")
	ErrInvalidResponse = errors.New("Invalid response")
	ErrInvalid         = errors.New("Invalid data")
	ErrNotNumber       = errors.New("Not a number")
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type response struct {
	Results      []item `json:"results"`
	Date         int    `json:"date"`
	TotalResults int    `json:"total_results"`
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return c
}

func getItems(url string) ([]item, error) {
	r, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("Can not read the Body: %s", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w, %s", err, msg)
	}
	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("%w,", ErrInvalid)
	}
	return resp.Results, nil
}

func getAll(apiUrl string) ([]item, error) {
	url := fmt.Sprintf("%s/todo", apiUrl)
	return getItems(url)
}

func getOne(apiUrl string, id int) (item, error) {
	url := fmt.Sprintf("%s/todo/%d", apiUrl, id)
	i, err := getItems(url)
	if err != nil {
		return item{}, err
	}
	return i[0], nil
}

func sendRequest(url, method, contentType string,
	expStatus int, body io.Reader) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", contentType)
	response, err := newClient().Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != expStatus {
		msg, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("Fail to read body: %w", err)
		}
		err = ErrInvalidResponse
		if response.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w, %s", err, msg)
	}
	return nil
}

func addItem(apiUrl, task string) error {
	u := fmt.Sprintf("%s/todo", apiUrl)

	item := struct {
		Task string
	}{
		Task: task,
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(item); err != nil {
		return err
	}
	return sendRequest(u, http.MethodPost, "application/json",
		http.StatusCreated, &buffer)
}

func completeItem(apiUrl string, id int) error {
	u := fmt.Sprintf("%s/todo/%d?complete", apiUrl, id)
	return sendRequest(u, http.MethodPatch, "", http.StatusOK, nil)
}

func deleteItem(apiUrl string, id int) error {
	u := fmt.Sprintf("%s/todo/%d", apiUrl, id)
	return sendRequest(u, http.MethodDelete, "", http.StatusNoContent, nil)
}
