package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/neinBit/ocg-games-service/internal/server/handler/api"
	"github.com/neinBit/ocg-games-service/internal/server/handler/rt"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/basenode"
)

func NewRouter() (http.Handler, error) {
	// Realtime
	baseNode := basenode.NewBaseNode()
	go baseNode.Run()

	err := baseNode.InitGames()
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)

	// Realtime
	router.Route("/rt", func(realtimeRouter chi.Router) {
		realtimeRouter.Get("/sse/game/{title}", func(w http.ResponseWriter, r *http.Request) {
			rt.HandleGameSSE(w, r, baseNode)
		})
		realtimeRouter.Get("/ws/game/{title}/room/{roomid}", func(w http.ResponseWriter, r *http.Request) {
			rt.HandleRoomWS(w, r, baseNode)
		})
	})

	// Json
	router.Route("/api", func(jsonRouter chi.Router) {
		jsonRouter.Post("/game/{title}/room", func(w http.ResponseWriter, r *http.Request) {
			api.HandleRoomPost(w, r, baseNode)
		})

		jsonRouter.Get("/game", func(w http.ResponseWriter, r *http.Request) {
			api.HandleGamesGet(w, r, baseNode)
		})
	})

	return router, nil
}
