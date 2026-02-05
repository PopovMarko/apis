package main

import (
	"log"
	"net/http"
	"sync"
)

func newMux(f string) http.Handler {
	m := http.NewServeMux()
	mutex := &sync.Mutex{}

	m.HandleFunc("/", rootHandler)

	handler := todoRouter(f, mutex)

	m.Handle("/todo", http.StripPrefix("/todo", handler))
	m.Handle("/todo/", http.StripPrefix("/todo/", handler))

	return m
}

func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content))
}

func replyJSONContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {
	body, err := resp.MarshallJSON()
	if err != nil {
		replyErrorContent(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func replyErrorContent(w http.ResponseWriter, r *http.Request, status int, err string) {
	log.Printf("%s, %s: Error: %d %s", r.URL, r.Method, status, err)
	http.Error(w, http.StatusText(status), status)

}
