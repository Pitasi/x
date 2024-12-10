package plausible

import (
	"net/http"
	"net/http/httputil"
)

var Proxy = &httputil.ReverseProxy{
	Director: func(r *http.Request) {
		r.Host = "plausible.anto.pt"
		r.URL.Scheme = "https"
		r.URL.Host = "plausible.anto.pt"
		if r.URL.Path == "/js/ps.js" {
			r.URL.Path = "/js/script.js"
		}
	},
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
		if r.Context().Err() != nil {
			return
		}
		http.Error(w, "proxy error", http.StatusBadGateway)
	},
}
