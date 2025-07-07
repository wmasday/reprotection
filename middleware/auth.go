package middleware

import (
	"net/http"
	"reprotection/config"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "session")

		userID, loggedIn := session.Values["user_id"]
		mustChange := session.Values["must_change_password"]

		// Redirect jika belum login
		if !loggedIn || userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Redirect ke /change-password jika wajib ganti password
		if mustChange == true && r.URL.Path != "/change-password" && r.URL.Path != "/change-password/submit" {
			http.Redirect(w, r, "/change-password", http.StatusSeeOther)
			return
		}

		// Lanjut ke handler berikutnya
		next(w, r)
	}
}
