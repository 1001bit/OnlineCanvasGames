package router

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/components"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/jsonapi"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/middleware"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/xmlapi"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/gamesservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/storageservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/userservice"
	"github.com/a-h/templ"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storageService *storageservice.Client, userService *userservice.Client, gamesService *gamesservice.Client) (http.Handler, error) {
	// Base
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
		jsonRouter.Post("/user/login", jsonapi.UserLoginHandler(userService))
		jsonRouter.Post("/user/register", jsonapi.UserRegisterHandler(userService))
		jsonRouter.Get("/logout", jsonapi.HandleLogout)

		// Secure routes
		jsonRouter.Group(func(jsonRouterSecure chi.Router) {
			jsonRouterSecure.Use(middleware.AuthJSON)

			jsonRouterSecure.Post("/game/{title}/room", gamesService.ProxyHandler())
		})
	})

	// XML
	router.Route("/", func(htmlRouter chi.Router) {
		htmlRouter.Get("/", templ.Handler(components.Home(gamesService)).ServeHTTP)
		htmlRouter.Get("/auth", xmlapi.HandleAuth)
		htmlRouter.Get("/profile/{name}", xmlapi.ProfileHandler(userService))

		// Secure routes
		htmlRouter.Group(func(htmlRouterSecure chi.Router) {
			htmlRouterSecure.Use(middleware.AuthHTML)

			htmlRouterSecure.Get("/game/{title}", xmlapi.HandleGameHub)
			htmlRouterSecure.Get("/game/{title}/room/{roomid}", xmlapi.HandleGameRoom)
		})

		// Non handled ones (404)
		htmlRouter.NotFound(templ.Handler(components.ErrorNotFound()).ServeHTTP)
	})

	return router, nil
}
