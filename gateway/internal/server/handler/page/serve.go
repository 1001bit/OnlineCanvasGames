package page

import (
	"html/template"
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/auth/claimscontext"
)

type NavigationData struct {
	Username string
}

type TemplateData struct {
	Data       any
	Navigation NavigationData
}

var templates = template.Must(template.ParseGlob("templates/**/*.html"))

func serveTemplate(w http.ResponseWriter, r *http.Request, file string, data any) {
	tmplData := TemplateData{
		Data: data,
	}

	tmplData.Navigation.Username, _ = claimscontext.GetUsername(r.Context())

	templates.ExecuteTemplate(w, file, tmplData)
}
