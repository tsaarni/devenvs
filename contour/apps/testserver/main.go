package main

import (
	"fmt"
	"net/http"
	"time"
)

func simpleResponse(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received:", req.RemoteAddr)
	fmt.Fprintf(w, "{response: 1}")
}

func slowStreamer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received:", req.RemoteAddr)

	cn, _ := w.(http.CloseNotifier)
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("not flushable")
	}

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	i := 0
	s := time.Duration(5)
	for i < 9999 {
		select {
		case <-cn.CloseNotify():
			fmt.Println("Client disconnected:", req.RemoteAddr)
			return
		default:
			fmt.Fprintf(w, "{response: %d, sleeping: %d}\n", i, s)
			flusher.Flush()
			i++
			time.Sleep(s * time.Second)
		}
	}
}

func main() {
	srv := &http.Server{
		Addr: ":8000",
		// ReadTimeout:       1 * time.Second,
		// WriteTimeout:      1 * time.Second,
		// IdleTimeout:       30 * time.Second,
		// ReadHeaderTimeout: 2 * time.Second,
		//Handler: http.HandlerFunc(slowStreamer),
		Handler: http.HandlerFunc(simpleResponse),
	}
	srv.SetKeepAlivesEnabled(true)

	fmt.Println("Listening for :8000")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
