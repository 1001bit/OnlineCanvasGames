package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	url *url.URL
}

func NewService(host, port string) (*Service, error) {
	service := &Service{}
	var err error

	service.url, err = url.Parse(fmt.Sprintf("http://%s:%s", host, port))
	return service, err
}

func (s *Service) Proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(s.url)
		proxy.ServeHTTP(w, r)
	}
}
