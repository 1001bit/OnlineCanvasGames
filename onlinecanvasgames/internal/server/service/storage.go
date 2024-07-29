package service

import "net/http"

type StorageService struct {
	service *Service
}

func NewStorageService(host, port string) (*StorageService, error) {
	service, err := NewService(host, port)
	if err != nil {
		return nil, err
	}
	return &StorageService{
		service: service,
	}, nil
}

func (s *StorageService) HandleStorage() http.Handler {
	return s.service.proxy()
}
