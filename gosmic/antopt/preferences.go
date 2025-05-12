package antopt

import (
	"html/template"
	"net/http"
	"slices"
)

func (ws Website) preferences(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r.Context())
	color := template.CSS(r.FormValue("color"))
	redirectTo := r.FormValue("redirectTo")
	if redirectTo == "" {
		redirectTo = "/"
	}

	user.Preferences.Background = ws.Colors[0]
	if slices.Contains(ws.Colors, color) {
		user.Preferences.Background = color
	}

	updatePrefs(w, user.Preferences)
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}
