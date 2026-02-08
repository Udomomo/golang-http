package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "Server is running!")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "Top page")
	})

	var h http.Handler = mux
	http.ListenAndServe("localhost:3000", h)
}
