package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func echoHandler(w http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed("POST", req) {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !isApplicationJson(req) {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func isMethodAllowed(method string, req *http.Request) bool {
	return req.Method == method
}

func isApplicationJson(req *http.Request) bool {
	return req.Header.Get("Content-Type") == "application/json"
}

func main() {
	mux := http.NewServeMux()

	// User型のjsonをPOSTリクエストで受け取り、バリデーションに通過すればそのまま返す。
	mux.HandleFunc("/echo", CorsMiddleware(echoHandler))

	// GETリクエストを受け取り、現在時刻を返す。
	mux.HandleFunc("/time", CorsMiddleware(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain")

		timezone, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, time.Now().In(timezone).Format("2006/01/02 15:04:05"))
	}))

	var h http.Handler = mux
	if err := http.ListenAndServe("localhost:3000", h); err != nil {
		fmt.Println("ListenAndServe error: ", err)
	}
}
