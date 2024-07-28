package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neinBit/ocg-storage-service/internal/router"
	"github.com/neinBit/ocg-storage-service/pkg/env"
)

func init() {
	env.InitEnv()
}

func main() {
	// start http server
	router, err := router.NewRouter()
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	addr := fmt.Sprintf(":%s", env.GetEnvVal("STORAGE_PORT"))

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
