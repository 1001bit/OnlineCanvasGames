package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/server/router"
	"github.com/1001bit/ocg-gateway-service/internal/server/service"
	"github.com/1001bit/ocg-gateway-service/pkg/env"
)

func main() {
	// services
	storageService, err := service.NewStorageService(env.GetEnvVal("STORAGE_HOST"), env.GetEnvVal("STORAGE_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	userService, err := service.NewUserService(env.GetEnvVal("USER_HOST"), env.GetEnvVal("USER_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	gamesService, err := service.NewGamesService(env.GetEnvVal("GAMES_HOST"), env.GetEnvVal("GAMES_PORT"))
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
