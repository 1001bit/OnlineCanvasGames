package page

import (
	"html/template"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
)

type HeaderData struct {
	Username string
	UserID   int
}

type TemplateData struct {
	Data   any
	Header HeaderData
}

var templates = template.Must(template.ParseGlob("web/templates/**/*.html"))

func serveTemplate(file string, data any, w http.ResponseWriter, r *http.Request) {
	tmplData := TemplateData{
		Data: data,
	}

	claims, ok := r.Context().Value(auth.ClaimsKey).(auth.Claims)
	if ok {
		tmplData.Header.UserID = claims.UserID
		tmplData.Header.Username = claims.Username
	}

	templates.ExecuteTemplate(w, file, tmplData)
}
