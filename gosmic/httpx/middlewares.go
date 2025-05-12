// Common HTTP middlewares.
package httpx

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func Logger(next http.Handler) http.Handler {
	// 2024/12/03 15:50:12 "GET http://localhost:8080/articles/help-your-code-reviewer HTTP/1.1" from [::1]:51394 - 200 16954B in 539.417µs
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info(r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"host", r.Host,
			"took", time.Since(start))
	})
}

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				slog.Error("panic", "err", rvr)
				fmt.Fprintln(os.Stderr, string(debug.Stack()))

				if r.Header.Get("Connection") != "Upgrade" {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gzipW := newCompressWriter(w)
		defer gzipW.Close()
		next.ServeHTTP(gzipW, r)
	})
}

type compressWriter struct {
	w  http.ResponseWriter
	gz *gzip.Writer

	wroteHeader bool
	compress    bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	gz := gzip.NewWriter(w)
	return &compressWriter{w: w, gz: gz}
}

func (c *compressWriter) Header() http.Header { return c.w.Header() }

func (c *compressWriter) WriteHeader(statusCode int) {
	if c.wroteHeader {
		c.w.WriteHeader(statusCode)
		return
	}

	c.wroteHeader = true
	defer c.w.WriteHeader(statusCode)

	contentType := c.Header().Get("Content-Type")
	contentEncoding := c.Header().Get("Content-Encoding")
	if contentEncoding != "" || !canCompress(contentType) {
		return
	}

	c.compress = true
	c.Header().Set("Content-Encoding", "gzip")
	c.Header().Add("Vary", "Accept-Encoding")
	c.Header().Del("Content-Length")
}

func (c *compressWriter) Write(b []byte) (int, error) {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}
	if c.compress {
		return c.gz.Write(b)
	} else {
		return c.w.Write(b)
	}
}

func (c *compressWriter) Flush() {
	if c.compress {
		c.gz.Flush()
	}

	if f, ok := c.w.(http.Flusher); ok {
		f.Flush()
	}
}

func (c *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := c.w.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("compress: http.Hijacker is unavailable on the writer")
}

func (c *compressWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := c.w.(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return errors.New("compress: http.Pusher is unavailable on the writer")
}

func (c *compressWriter) Close() error {
	if c.compress {
		c.gz.Close()
	}

	if w, ok := c.w.(io.Closer); ok {
		return w.Close()
	}

	return errors.New("compress: io.WriteCloser is unavailable on the writer")
}

func (c *compressWriter) Unwrap() http.ResponseWriter {
	return c.w
}

var _ http.ResponseWriter = &compressWriter{}
var _ http.Flusher = &compressWriter{}

func canCompress(contentType string) bool {
	ct := strings.SplitN(contentType, ";", 2)[0]
	switch ct {
	case "text/html", "text/css", "text/plain":
		return true
	}
	return false
}

var httpReqs = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "HTTP requests processed, partitioned by handler.",
}, []string{"host", "handler"})

func MetricsInc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpReqs.WithLabelValues(r.Host, r.URL.String()).Inc()
		h.ServeHTTP(w, r)
	})
}
