package config

import "github.com/gorilla/sessions"

var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte("your-secret-key"))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}
}
