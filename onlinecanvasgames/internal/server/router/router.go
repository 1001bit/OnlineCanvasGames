package router

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/rt"
	"github.com/1001bit/OnlineCanvasGames/internal/server/middleware"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/basenode"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storageService *Service) (http.Handler, error) {
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)
	router.Use(middleware.InjectJWTClaims)

	// Storage
	router.Handle("/static/*", storageService.Proxy())
	router.Get("/favicon.ico", storageService.Proxy())
	router.Handle("/js/*", storageService.Proxy())
	router.Handle("/image/*", storageService.Proxy())

	// Realtime
	baseNode := basenode.NewBaseNode()
	go baseNode.Run()

	err := baseNode.InitGames()
	if err != nil {
		return nil, err
	}

	// RT Secure
	router.Route("/rt", func(rs chi.Router) {
		rs.Use(middleware.AuthJSON)

		rs.Get("/sse/game/{gameid}", func(w http.ResponseWriter, r *http.Request) {
			rt.HandleGameSSE(w, r, baseNode)
		})
		rs.Get("/ws/game/{gameid}/room/{roomid}", func(w http.ResponseWriter, r *http.Request) {
			rt.HandleRoomWS(w, r, baseNode)
		})
	})

	// API
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.TypeJSON)

		// Post
		r.Post("/user", api.HandleUserPost)
		r.Post("/game/{gameid}/room", func(w http.ResponseWriter, r *http.Request) {
			api.HandleRoomPost(w, r, baseNode)
		})
	})

	// HTML Pages
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.TypeHTML)

		// Get
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			page.HandleHome(w, r, baseNode)
		})
		r.Get("/auth", page.HandleAuth)
		r.Get("/profile/{id}", page.HandleProfile)
		r.Get("/logout", page.HandleLogout)

		// Secure
		r.Group(func(rs chi.Router) {
			rs.Use(middleware.AuthHTML)

			rs.Get("/game/{gameid}", page.HandleGameHub)
			rs.Get("/game/{gameid}/room/{roomid}", page.HandleGameRoom)
		})

		r.Get("/*", page.HandleNotFound)
	})

	return router, nil
}