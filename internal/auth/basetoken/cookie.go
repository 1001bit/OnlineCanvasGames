package basetoken

import (
	"net/http"
	"time"
)

func NewCookie(tokenStr string, name string, exp time.Duration) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   int(exp.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
}
