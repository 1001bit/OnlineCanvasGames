package page

import (
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access",
		MaxAge: 0,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh",
		MaxAge: 0,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
