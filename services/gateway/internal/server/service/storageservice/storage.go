package storageservice

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service"
)

type StorageService struct {
	service *service.Rest
}

func New(host, port string) (*StorageService, error) {
	service, err := service.NewRestService(host, port)
	if err != nil {
		return nil, err
	}
	return &StorageService{
		service: service,
	}, nil
}

func (s *StorageService) ProxyHandler() http.Handler {
	return s.service.Proxy()
}
