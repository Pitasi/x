package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"anto.pt/x/wlog"
)

func AuthenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // simulate fetching of user session

		// enrich your widelog event
		wlog.Group(r.Context(), "user", "id", 1234, "name", "foo")

		h.ServeHTTP(w, r)
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond) // simulate request time
	wlog.Set(r.Context(), "cart_status", "empty")
	w.Write([]byte("OK!"))
}

func handleForce(w http.ResponseWriter, r *http.Request) {
	wlog.MustLog(r.Context()) // force this request to be logged
	w.Write([]byte("OK!"))
}

func handleError(w http.ResponseWriter, r *http.Request) {
	wlog.Set(r.Context(), "err", "something bad happened") // set an error field
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("err!"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /", AuthenticationMiddleware(http.HandlerFunc(handle)))
	mux.Handle("GET /force", http.HandlerFunc(handleForce))
	mux.Handle("GET /error", http.HandlerFunc(handleError))

	handler := wlog.Middleware(mux, slog.Default(), wlog.NewSampler(0.5))

	fmt.Println("listening at :9090")
	http.ListenAndServe(":9090", handler)
}
