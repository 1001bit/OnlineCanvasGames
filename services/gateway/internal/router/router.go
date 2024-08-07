package router

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/handler/api"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/handler/page"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/middleware"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/gamesservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/storageservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/userservice"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storageService *storageservice.Client, userService *userservice.Client, gamesService *gamesservice.Client) (http.Handler, error) {
	// Router
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)
	router.Use(middleware.ClaimsContext)

	// Storage
	router.Group(func(storageRouter chi.Router) {
		storageRouter.Handle("/css/*", storageService.ProxyHandler())
		storageRouter.Handle("/favicon.ico", storageService.ProxyHandler())
		storageRouter.Handle("/js/*", storageService.ProxyHandler())
		storageRouter.Handle("/image/*", storageService.ProxyHandler())
	})

	// Realtime
	router.Route("/rt", func(realtimeRouter chi.Router) {
		realtimeRouter.Use(middleware.AuthJSON)

		realtimeRouter.Get("/sse/game/{title}", gamesService.ProxyHandler())
		realtimeRouter.Get("/ws/game/{title}/room/{roomid}", gamesService.RoomProxyHandler())
	})

	// JSON
	router.Route("/api", func(jsonRouter chi.Router) {
		// Login
		jsonRouter.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
			api.HandleUserLogin(w, r, userService)
		})
		// Register
		jsonRouter.Post("/user/register", func(w http.ResponseWriter, r *http.Request) {
			api.HandleUserRegister(w, r, userService)
		})

		// Rooms
		jsonRouter.Group(func(jsonRouterSecure chi.Router) {
			jsonRouterSecure.Use(middleware.AuthJSON)

			jsonRouterSecure.Post("/game/{title}/room", gamesService.ProxyHandler())
		})
	})

	// HTML Pages
	router.Route("/", func(htmlRouter chi.Router) {
		// Home
		htmlRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			page.HandleHome(w, r, gamesService)
		})
		// Auth
		htmlRouter.Get("/auth", page.HandleAuth)
		// Profile
		htmlRouter.Get("/profile/{name}", func(w http.ResponseWriter, r *http.Request) {
			page.HandleProfile(w, r, userService)
		})
		// Logout
		htmlRouter.Get("/logout", page.HandleLogout)

		// Games
		htmlRouter.Group(func(htmlRouterSecure chi.Router) {
			htmlRouterSecure.Use(middleware.AuthHTML)

			htmlRouterSecure.Get("/game/{title}", page.HandleGameHub)
			htmlRouterSecure.Get("/game/{title}/room/{roomid}", page.HandleGameRoom)
		})

		// Other
		htmlRouter.NotFound(page.HandleNotFound)
	})

	return router, nil
}
