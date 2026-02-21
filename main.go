package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Counter struct {
	mu    sync.Mutex
	count int
}

// ハンドラ関数からCounterを操作できるようレシーバ関数にしている。
func (c *Counter) handleGet(w http.ResponseWriter, req *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(w, c.count)
}

func (c *Counter) handleInc(w http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed("POST", req) {
		w.Header().Add("Content-type", "text/plain")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(w, c.count)
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

	// カウンタの実装
	c := &Counter{
		count: 0,
	}
	mux.HandleFunc("/count", CorsMiddleware(c.handleGet))
	mux.HandleFunc("/inc", CorsMiddleware(c.handleInc))

	// sleepリクエスト (graceful shutdownテスト用)
	mux.HandleFunc("/sleep", CorsMiddleware(func(w http.ResponseWriter, req *http.Request) {
		for i := 1; i <= 10; i++ {
			fmt.Println("sleep", i, "秒経過")
			time.Sleep(1 * time.Second)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "slept 10 seconds")
	}))

	var h http.Handler = mux
	srv := &http.Server{
		Addr:              ":3000",
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		fmt.Println("Starting server on localhost:3000...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("ListenAndServe error: ", err)
		}
	}()

	// Ctrl+Cでgraceful Shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("HTTP server Shutdown: %v", err)
	}

	fmt.Println("Gracefully exiting...")
}
