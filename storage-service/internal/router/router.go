package router

import (
	"net/http"

	"github.com/1001bit/ocg-storage-service/internal/handler"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter() (http.Handler, error) {
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.RedirectSlashes)

	// Image
	router.Handle("/image/*", http.StripPrefix("/image", http.HandlerFunc(handler.HandleImage)))

	// Everything else
	router.Handle("/*", handler.StaticHandler("/"))

	return router, nil
}
