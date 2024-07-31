package handler

import (
	"net/http"
	"path"
)

const staticPath = "./static"

func StaticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := path.Join(staticPath, dir)
		fileServer := http.FileServer(http.Dir(path))
		fileServer.ServeHTTP(w, r)
	}
}
