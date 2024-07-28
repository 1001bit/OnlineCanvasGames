package main

import (
	"log"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/basetoken"
	"github.com/1001bit/OnlineCanvasGames/internal/database"
	"github.com/1001bit/OnlineCanvasGames/internal/server/router"
	"github.com/1001bit/OnlineCanvasGames/pkg/env"
)

func init() {
	env.InitEnv()
	basetoken.InitJWTSecret()
}

func main() {
	// start database
	err := database.Start()
	if err != nil {
		log.Fatal("err starting database:", err)
	}
	defer database.DB.Close()

	// start http server
	router, err := router.NewRouter()
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	addr := env.GetEnvVal("ADDR")

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
