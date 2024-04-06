package router

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/sse"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/storage"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/ws"
	"github.com/1001bit/OnlineCanvasGames/internal/server/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter() (http.Handler, error) {
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)

	// Storage
	router.Handle("/static/*", http.StripPrefix("/static", http.HandlerFunc(storage.HandleStatic)))
	router.Get("/favicon.ico", storage.HandleStatic)
	router.Handle("/image/*", http.StripPrefix("/image", http.HandlerFunc(storage.HandleImage)))

	// Websockets
	gamesWS := ws.NewGamesWS()
	go gamesWS.Run()
	// WS Secure
	router.Route("/ws", func(rs chi.Router) {
		rs.Use(middleware.AuthJSON)

		rs.HandleFunc("/game/{id}/{roomid}", gamesWS.HandleWS)
	})

	// Server-Sent Events
	gamesSSE, err := sse.NewGamesSSE()
	if err != nil {
		return nil, err
	}
	go gamesSSE.Run()
	// SSE Secure
	router.Route("/sse", func(rs chi.Router) {
		rs.Use(middleware.AuthJSON)

		rs.HandleFunc("/game/{id}", gamesSSE.HandleEvent)
	})

	// API
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.TypeJSON)

		// Post
		r.Post("/user", api.HandleUserPost)
	})

	// HTML Pages
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.TypeHTML)

		// Get
		r.Get("/", page.HandleHome)
		r.Get("/auth", page.HandleAuth)
		r.Get("/profile/{id}", page.HandleProfile)
		// Secure
		r.Group(func(rs chi.Router) {
			rs.Use(middleware.AuthHTML)

			rs.Get("/game/{gameid}/hub", page.HandleGameHub)
			rs.Get("/game/{gameid}/room/{roomid}", page.HandleGameRoom)
		})

		r.Get("/*", page.HandleNotFound)
	})

	return router, nil
}
