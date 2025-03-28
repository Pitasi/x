package main

import (
	"log"
	"net/http"
)

func colophon(mux *http.ServeMux) {
	mux.HandleFunc("GET /colophon", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "colophon.html", struct {
			Site   Site
			Active string
			User   *User
		}{
			Site:   site,
			User:   user,
			Active: "/colophon",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})
}
