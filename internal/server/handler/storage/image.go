package storage

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
)

const imageStoragePath = "./web/image"

var imageFormats = [...]string{"png", "jpg", "jpeg"}

func HandleImage(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if strings.Contains(path, "..") {
		api.ServeJSONMessage("invalid path", http.StatusBadRequest, w)
		return
	}

	// check if file.png, file.jpg, etc. exists and use it
	for _, format := range imageFormats {
		possiblePath := filepath.Join(imageStoragePath, path+"."+format)

		if _, err := os.Stat(possiblePath); err == nil {
			http.ServeFile(w, r, possiblePath)
			return
		}
	}

	fileServer := http.FileServer(http.Dir(imageStoragePath))
	fileServer.ServeHTTP(w, r)
}
