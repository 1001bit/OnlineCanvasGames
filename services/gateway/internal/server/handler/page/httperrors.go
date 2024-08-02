package page

import "net/http"

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	serveTemplate(w, r, "notFound.html", nil)
}

func HandleServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	serveTemplate(w, r, "serverError.html", nil)
}
