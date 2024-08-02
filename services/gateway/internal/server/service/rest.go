package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/message"
)

type Rest struct {
	url *url.URL
}

func NewRestService(host, port string) (*Rest, error) {
	url, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))
	return &Rest{
		url: url,
	}, err
}

func (s *Rest) Proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(s.url)
		proxy.ServeHTTP(w, r)
	}
}

func (s *Rest) Request(ctx context.Context, method string, path string, requestBody io.Reader) (*message.JSON, error) {
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
