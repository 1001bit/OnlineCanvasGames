package router

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/components"
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

	// API
	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Post("/login", xmlapi.LoginHandler(userService))
		apiRouter.Post("/register", xmlapi.RegisterHandler(userService))
		apiRouter.Get("/logout", xmlapi.HandleLogout)

		// Secure
		apiRouter.Group(func(apiRouterSecure chi.Router) {
			apiRouterSecure.Use(middleware.AuthJSON)

			apiRouterSecure.Post("/game/{title}/room", gamesService.ProxyHandler())
		})
	})

	// Pages
	router.Get("/", templ.Handler(components.Home(gamesService)).ServeHTTP)
	router.Get("/auth", xmlapi.HandleAuthPage)
	router.Get("/profile/{name}", xmlapi.ProfileHandler(userService))

	// Secure
	router.Group(func(routerSecure chi.Router) {
		routerSecure.Use(middleware.AuthHTML)

		routerSecure.Get("/game/{title}", xmlapi.HandleGameHub)
		routerSecure.Get("/game/{title}/room/{roomid}", xmlapi.HandleGameRoom)
	})

	// Non handled ones (404)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		components.ErrorNotFound().Render(r.Context(), w)
	})

	return router, nil
}
