# wlog

Inspired by https://loggingsucks.com/, this is my minimalistic approach for
having wide logs in Go using `slog`.

See [example/main.go](./example/main.go) for a full example using it.

## Installation

```sh
go get anto.pt/x/wlog
```

## Usage

```go
mux := http.NewServeMux()

// wrap your http handler with the wlog middleware
handler := wlog.Middleware(mux, slog.Default(), wlog.NewSampler(0.5))
http.ListenAndServe(":9090", handler)

// ...

// use wlog in you http handlers or middlewares
func(w http.ResponseWriter, r *http.Request) {
    // enrich your wlog event
    wlog.Group(r.Context(), "user", "id", 1234, "name", "foo")
    w.Write([]byte("hello!"))
}

// get your log:
// 2026/02/06 16:58:12 INFO request user.id=1234 user.name=foo took=202.235667ms
```

See [example/main.go](./example/main.go) for a full example using it.
