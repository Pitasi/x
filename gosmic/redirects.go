package main

import (
	"net/http"
)

func redirects(mux *http.ServeMux) {
	var maps = map[string]string{
		"/articles/": "/",
		"/gomad24/":  "https://docs.google.com/presentation/d/1Z8GZRU9H_OEzigLn5riHF6_Io7zwrtnXzgMXEABm0lc",
	}

	for from, to := range maps {
		mux.HandleFunc(from, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, to, http.StatusFound)
		})
	}
}
