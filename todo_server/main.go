package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	host := flag.String("h", "localhost", "Server host")
	port := flag.Int("p", 8080, "Server port")
	todoFile := flag.String("f", "todo_server.json", "File name to store")
	flag.Parse()

	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", *host, *port),
		Handler:      newMux(*todoFile),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server started at %s, port %d\nUsing file %s", *host, *port, *todoFile)
	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "Fail to start server: %s", err)
		os.Exit(1)
	}
}
