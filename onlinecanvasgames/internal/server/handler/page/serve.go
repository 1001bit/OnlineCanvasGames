package page

import (
	"html/template"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
)

type NavigationData struct {
	Username string
	UserID   int
}

type TemplateData struct {
	Data       any
	Navigation NavigationData
}

var templates = template.Must(template.ParseGlob("web/templates/**/*.html"))

func serveTemplate(file string, data any, w http.ResponseWriter, r *http.Request) {
	tmplData := TemplateData{
		Data: data,
	}

	claims, err := auth.GetJwtClaimsFromContext(r.Context())
	if err == nil {
		tmplData.Navigation.UserID = claims.UserID
		tmplData.Navigation.Username = claims.Username
	}

	templates.ExecuteTemplate(w, file, tmplData)
}
