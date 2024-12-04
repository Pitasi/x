package main

import (
	"html/template"
	"net/http"
)

func prefs(mux *http.ServeMux) {
	mux.HandleFunc("POST /submit-color", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value("user").(*User)
		color := template.CSS(r.FormValue("color"))
		redirectTo := r.FormValue("redirectTo")
		if redirectTo == "" {
			redirectTo = "/"
		}

		user.Preferences.Background = site.Colors[0]
		for _, c := range site.Colors {
			if c == color {
				user.Preferences.Background = c
				break
			}
		}

		updatePrefs(w, user.Preferences)
		http.Redirect(w, r, r.FormValue("redirectTo"), http.StatusSeeOther)
	})
}
