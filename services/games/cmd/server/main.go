package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/database"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/gamemodel"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/router"
	"github.com/1001bit/overenv"
)

func main() {
	// start database
	db, err := database.NewFromEnv()
	if err != nil {
		log.Fatal("err starting database:", err)
	}
	defer db.Close()

	// gamestore
	gameStore := gamemodel.NewGameStore(db)

	// router
	router, err := router.NewRouter(gameStore)
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	// start http server
	addr := fmt.Sprintf(":%s", overenv.Get("PORT"))
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
