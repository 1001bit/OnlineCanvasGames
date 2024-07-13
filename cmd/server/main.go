package main

import (
	"log"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/basetoken"
	"github.com/1001bit/OnlineCanvasGames/internal/database"
	"github.com/1001bit/OnlineCanvasGames/internal/env"
	"github.com/1001bit/OnlineCanvasGames/internal/server/router"
	"github.com/1001bit/OnlineCanvasGames/internal/tscompiler"
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

	// compile typescript
	err = tscompiler.CompileTypeScript()
	if err != nil {
		log.Fatal("err compiling typescript:", err)
	}

	// start http server
	router, err := router.NewRouter()
	if err != nil {
		log.Fatal("err creating router:", err)
	}

	addr := env.GetEnvVal("ADDR")

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
