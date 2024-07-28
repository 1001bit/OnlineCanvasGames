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

	// Static
	router.Handle("/static/*", handler.StorageHandler("/"))
	router.Get("/favicon.ico", handler.StorageHandler("/static"))

	// JS
	router.Handle("/js/*", handler.StorageHandler("/"))

	// Image
	router.Handle("/image/*", http.StripPrefix("/image", http.HandlerFunc(handler.HandleImage)))

	return router, nil
}
