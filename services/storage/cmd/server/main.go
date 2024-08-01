package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/storage/internal/router"
	"github.com/1001bit/onlinecanvasgames/services/storage/pkg/env"
)

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
