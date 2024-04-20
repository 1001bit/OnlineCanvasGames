package router

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/storage"
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
	router.Handle("/gamescript/*", http.StripPrefix("/gamescript", http.HandlerFunc(storage.HandleGamescript)))

	// Realtime
	rt := realtime.NewRealtime()
	go rt.Run()

	err := rt.InitGames()
	if err != nil {
		return nil, err
	}

	// RT Secure
	router.Route("/rt", func(rs chi.Router) {
		rs.Use(middleware.AuthJSON)

		rs.Get("/sse/game/{gameid}", rt.HandleGameSSE)
		rs.Get("/ws/game/{gameid}/room/{roomid}", rt.HandleRoomWS)
	})

	// API
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.TypeJSON)

		// Post
		r.Post("/user", api.HandleUserPost)
		r.Post("/game/{gameid}/room", func(w http.ResponseWriter, r *http.Request) {
			api.HandleRoomPost(w, r, rt)
		})

		// Get
		r.Get("/game/{gameid}/room", func(w http.ResponseWriter, r *http.Request) {
			api.HandleRandomRoomGet(w, r, rt)
		})
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

			rs.Get("/game/{gameid}", page.HandleGame)
			rs.Get("/game/{gameid}/room/{roomid}", page.HandleGameRoom)
		})

		r.Get("/*", page.HandleNotFound)
	})

	return router, nil
}
