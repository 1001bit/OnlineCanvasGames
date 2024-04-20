package storage

import "net/http"

const gamescriptStoragePath = "./web/gamescript"

func HandleGamescript(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.Dir(gamescriptStoragePath))
	fileServer.ServeHTTP(w, r)
}
