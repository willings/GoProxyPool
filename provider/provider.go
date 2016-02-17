package provider

import "net/http"

type ProxyItem struct {
	Host      string `json:"Host"`
	Port      int    `json:"Port"`
	Type      int    `json:"Type"`
	Anonymous int    `json:"Anonymous"`
}

type ProxyProvider interface {
	Load() ([]*ProxyItem, error)

	SetClient(*http.Client)
}
