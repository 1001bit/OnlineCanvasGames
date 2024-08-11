package jsonapi

import (
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access",
		Path:   "/",
		MaxAge: 0,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh",
		Path:   "/",
		MaxAge: 0,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
