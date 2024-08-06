package storageservice

import (
	"net/http"

	service "github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client"
)

type Client struct {
	service *service.RestClient
}

func NewClient(host, port string) (*Client, error) {
	service, err := service.NewRestClient(host, port)
	if err != nil {
		return nil, err
	}
	return &Client{
		service: service,
	}, nil
}

func (s *Client) ProxyHandler() http.Handler {
	return s.service.Proxy()
}
