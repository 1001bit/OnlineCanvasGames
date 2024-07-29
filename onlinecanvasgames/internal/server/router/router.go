package router

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/rt"
	"github.com/1001bit/OnlineCanvasGames/internal/server/middleware"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/basenode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/service"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storageService *service.StorageService, userService *service.UserService) (http.Handler, error) {
	// Realtime
	baseNode := basenode.NewBaseNode()
	go baseNode.Run()

	err := baseNode.InitGames()
	if err != nil {
		return nil, err
	}

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

		realtimeRouter.Get("/sse/game/{gameid}", func(w http.ResponseWriter, r *http.Request) {
			rt.HandleGameSSE(w, r, baseNode)
		})
		realtimeRouter.Get("/ws/game/{gameid}/room/{roomid}", func(w http.ResponseWriter, r *http.Request) {
			rt.HandleRoomWS(w, r, baseNode)
		})
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

			jsonRouterSecure.Post("/game/{gameid}/room", func(w http.ResponseWriter, r *http.Request) {
				api.HandleRoomPost(w, r, baseNode)
			})
		})
	})

	// HTML Pages
	router.Route("/", func(htmlRouter chi.Router) {
		htmlRouter.Use(middleware.TypeHTML)

		// Home
		htmlRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			page.HandleHome(w, r, baseNode)
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

			htmlRouterSecure.Get("/game/{gameid}", page.HandleGameHub)
			htmlRouterSecure.Get("/game/{gameid}/room/{roomid}", page.HandleGameRoom)
		})

		// Other
		htmlRouter.Get("/*", page.HandleNotFound)
	})

	return router, nil
}
