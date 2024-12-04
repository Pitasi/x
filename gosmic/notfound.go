package main

import (
	"log/slog"
	"net/http"
)

func notfound(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		err := templates.ExecuteTemplate(w, "notfound.html", struct {
			Site   Site
			User   *User
			Active string
		}{
			Site:   site,
			User:   user,
			Active: "",
		})
		if err != nil {
			slog.Error("error rendering notfound.html", "err", err)
		}
	})
}
