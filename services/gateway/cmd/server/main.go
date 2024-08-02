package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/router"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/gamesservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/storageservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/userservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/env"
)

func main() {
	// services
	storageService, err := storageservice.New(env.GetEnvVal("STORAGE_HOST"), env.GetEnvVal("STORAGE_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	userService, err := userservice.New(env.GetEnvVal("USER_HOST"), env.GetEnvVal("USER_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	gamesService, err := gamesservice.New(env.GetEnvVal("GAMES_HOST"), env.GetEnvVal("GAMES_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	// router
	router, err := router.NewRouter(storageService, userService, gamesService)
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	// start http server
	addr := fmt.Sprintf(":%s", env.GetEnvVal("PORT"))
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
