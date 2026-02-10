package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/time", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain")

		timezone, err := time.LoadLocation("Asia/Tokyo")
    if err != nil {
        panic(err)
    }
		fmt.Fprint(w, time.Now().In(timezone).Format("2006/01/02 15:04:05"))
	})

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
