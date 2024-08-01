package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/ocg-user-service/internal/database"
	"github.com/1001bit/ocg-user-service/internal/server/router"
	"github.com/1001bit/ocg-user-service/pkg/env"
)

func main() {
	// start database
	err := database.Start()
	if err != nil {
		log.Fatal("err starting database:", err)
	}
	defer database.DB.Close()

	// router
	router, err := router.NewRouter()
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	// start http server
	addr := fmt.Sprintf(":%s", env.GetEnvVal("PORT"))
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
