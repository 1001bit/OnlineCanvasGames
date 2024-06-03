package storage

import "net/http"

const gameAssetsStoragePath = "./web/gameassets"

func HandleGameAssets(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.Dir(gameAssetsStoragePath))
	fileServer.ServeHTTP(w, r)
}
