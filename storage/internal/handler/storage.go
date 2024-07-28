package handler

import (
	"net/http"
	"path"
)

const storagePath = "./storage"

func StorageHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := path.Join(storagePath, dir)
		fileServer := http.FileServer(http.Dir(path))
		fileServer.ServeHTTP(w, r)
	}
}
