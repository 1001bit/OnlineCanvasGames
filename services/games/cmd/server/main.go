package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/database"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/router"
	"github.com/1001bit/onlinecanvasgames/services/games/pkg/env"
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