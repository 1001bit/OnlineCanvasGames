package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/neinBit/ocg-storage-service/internal/handler"
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
