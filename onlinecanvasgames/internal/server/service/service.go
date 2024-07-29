package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
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

func (s *Service) proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(s.url)
		proxy.ServeHTTP(w, r)
	}
}

func (s *Service) request(ctx context.Context, method string, path string, requestBody io.Reader) (*message.JSON, error) {
	req, err := http.NewRequest(method, s.url.JoinPath(path).String(), requestBody)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	msg := &message.JSON{}
	err = json.NewDecoder(resp.Body).Decode(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
