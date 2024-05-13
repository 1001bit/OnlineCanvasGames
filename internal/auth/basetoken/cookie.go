package basetoken

import (
	"net/http"
	"time"
)

const secure = false

func NewCookie(tokenStr string, name string, exp time.Duration) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   int(exp.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
}
