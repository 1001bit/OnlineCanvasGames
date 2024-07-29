package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/neinBit/ocg-user-service/internal/server/handler"
	"github.com/neinBit/ocg-user-service/internal/server/middleware"
)

func NewRouter() (http.Handler, error) {
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)
	router.Use(middleware.TypeJSON)

	// Post
	router.Post("/", handler.HandleUserPost)
	// Get
	router.Get("/{id}", handler.HandleUserGet)

	return router, nil
}
