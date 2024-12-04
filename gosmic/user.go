package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

type User struct {
	Preferences *UserPreferences
}

type UserPreferences struct {
	Background template.CSS
}

var defaultPrefs = UserPreferences{}
var defaultPrefsJSON, _ = json.Marshal(defaultPrefs)

func encodeCookie(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeCookie(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func defaultCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "prefs",
		Value:    encodeCookie(defaultPrefsJSON),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
	}
}

func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &User{
			Preferences: &UserPreferences{},
		}

		var value []byte
		cookie, err := r.Cookie("prefs")
		if err == nil {
			value, _ = decodeCookie(cookie.Value)
		}
		if value == nil {
			value = defaultPrefsJSON
		}

		err = json.Unmarshal(value, user.Preferences)
		if err != nil {
			// do nothing
		}

		// add User to context
		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		next.ServeHTTP(w, r)
	})
}

func updatePrefs(w http.ResponseWriter, prefs *UserPreferences) {
	prefsJson, err := json.Marshal(prefs)
	if err != nil {
		log.Println("marhsalling prefs cookie", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "prefs",
		Value:    encodeCookie(prefsJson),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
	})
}
