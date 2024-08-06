package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/router"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/gamesservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/storageservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/userservice"
	"github.com/1001bit/overenv"
)

func main() {
	// services
	storageService, err := storageservice.NewClient(overenv.Get("STORAGE_HOST"), overenv.Get("STORAGE_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	userService, err := userservice.NewClient(overenv.Get("USER_HOST"), overenv.Get("USER_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	gamesService, err := gamesservice.NewClient(overenv.Get("GAMES_HOST"), overenv.Get("GAMES_PORT"))
	if err != nil {
		log.Fatal("err getting service url:", err)
	}

	// router
	router, err := router.NewRouter(storageService, userService, gamesService)
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	// start http server
	addr := fmt.Sprintf(":%s", overenv.Get("PORT"))
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
