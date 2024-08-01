package storage

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service"
)

func Handler(service *service.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.ProxyHandler().ServeHTTP(w, r)
	}
}