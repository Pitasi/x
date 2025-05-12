package antopt

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

		_ = json.Unmarshal(value, user.Preferences)

		// add User to context
		r = r.WithContext(context.WithValue(r.Context(), userKey{}, user))

		next.ServeHTTP(w, r)
	})
}

type userKey struct{}

func GetUser(ctx context.Context) *User {
	user, _ := ctx.Value(userKey{}).(*User)
	return user
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
