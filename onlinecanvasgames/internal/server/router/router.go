package router

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
	"github.com/1001bit/OnlineCanvasGames/internal/server/middleware"
	"github.com/1001bit/OnlineCanvasGames/internal/server/service"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storageService *service.StorageService, userService *service.UserService, gamesService *service.GamesService) (http.Handler, error) {
	// Router
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)
	router.Use(middleware.ClaimsContext)

	// Storage
	router.Group(func(storageRouter chi.Router) {
		storageRouter.Handle("/static/*", storageService.HandleStorage())
		storageRouter.Get("/favicon.ico", storageService.HandleStorage().ServeHTTP)
		storageRouter.Handle("/js/*", storageService.HandleStorage())
		storageRouter.Handle("/image/*", storageService.HandleStorage())
	})

	// Realtime
	router.Route("/rt", func(realtimeRouter chi.Router) {
		realtimeRouter.Use(middleware.AuthJSON)

		realtimeRouter.Get("/sse/game/{gameid}", gamesService.HandleGameHubSSE())
		realtimeRouter.Get("/ws/game/{gameid}/room/{roomid}", gamesService.HandleRoomWS())
	})

	// JSON
	router.Route("/api", func(jsonRouter chi.Router) {
		jsonRouter.Use(middleware.TypeJSON)

		// Users
		jsonRouter.Post("/user", func(w http.ResponseWriter, r *http.Request) {
			api.HandleUserPost(w, r, userService)
		})

		// Rooms
		jsonRouter.Group(func(jsonRouterSecure chi.Router) {
			jsonRouterSecure.Use(middleware.AuthJSON)

			jsonRouterSecure.Post("/game/{gameid}/room", gamesService.HandleRoomPost())
		})
	})

	// HTML Pages
	router.Route("/", func(htmlRouter chi.Router) {
		htmlRouter.Use(middleware.TypeHTML)

		// Home
		htmlRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			page.HandleHome(w, r, gamesService)
		})
		// Auth
		htmlRouter.Get("/auth", page.HandleAuth)
		// Profile
		htmlRouter.Get("/profile/{id}", func(w http.ResponseWriter, r *http.Request) {
			page.HandleProfile(w, r, userService)
		})
		// Logout
		htmlRouter.Get("/logout", page.HandleLogout)

		// Games
		htmlRouter.Group(func(htmlRouterSecure chi.Router) {
			htmlRouterSecure.Use(middleware.AuthHTML)

			htmlRouterSecure.Get("/game/{gameid}", func(w http.ResponseWriter, r *http.Request) {
				page.HandleGameHub(w, r, gamesService)
			})
			htmlRouterSecure.Get("/game/{gameid}/room/{roomid}", func(w http.ResponseWriter, r *http.Request) {
				page.HandleGameRoom(w, r, gamesService)
			})
		})

		// Other
		htmlRouter.Get("/*", page.HandleNotFound)
	})

	return router, nil
}
