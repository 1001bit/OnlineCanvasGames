package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/database"
	"github.com/1001bit/OnlineCanvasGames/internal/server/router"
	"github.com/1001bit/OnlineCanvasGames/internal/server/service"
	"github.com/1001bit/OnlineCanvasGames/pkg/env"
)

func main() {
	// start database
	err := database.Start()
	if err != nil {
		log.Fatal("err starting database:", err)
	}
	defer database.DB.Close()

	// services
	storageService, err := service.NewStorageService(env.GetEnvVal("STORAGE_HOST"), env.GetEnvVal("STORAGE_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	userService, err := service.NewUserService(env.GetEnvVal("USER_HOST"), env.GetEnvVal("USER_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	// router
	router, err := router.NewRouter(storageService, userService)
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	// start http server
	addr := fmt.Sprintf(":%s", env.GetEnvVal("PORT"))
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
