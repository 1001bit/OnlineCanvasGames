package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const imageDir = "image"

var imageFormats = [...]string{"png", "jpg", "jpeg"}

func HandleImage(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if strings.Contains(path, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imagePath := filepath.Join(staticPath, imageDir)

	// check if file.png, file.jpg, etc. exists and use it
	for _, format := range imageFormats {
		possiblePath := filepath.Join(imagePath, path+"."+format)

		if _, err := os.Stat(possiblePath); err == nil {
			http.ServeFile(w, r, possiblePath)
			return
		}
	}

	fileServer := http.FileServer(http.Dir(imagePath))
	fileServer.ServeHTTP(w, r)
}
