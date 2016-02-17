package proxypool

import (
	"errors"
	"fmt"
	"github.com/willings/proxypool/provider"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Validator         ValidateFunc
	ValidateInterval  int
	maxValidateThread int

	Filter Filter

	Provider       provider.ProxyProvider
	ReloadInterval int

	ProxyStrategy ProxyStrategy
	AllowDirect   bool
}

func DefaultConfig() *Config {
	return &Config{
		Validator: func(p ProxyInfo) (*ProxyQuality, error) {
			urlStr := fmt.Sprintf("http://%s:%d", p.Host, p.Port)
			url, err := url.Parse(urlStr)
			if err != nil {
				return nil, err
			}
			proxy := http.ProxyURL(url)

			client := &http.Client{
				Transport: &http.Transport{
					Proxy: proxy,
				},
				Timeout: time.Duration(10 * time.Second),
			}

			startTime := time.Now()
			resp, err := client.Get("http://search.yahoo.co.jp/search?p=%E6%97%A5%E6%9C%AC")

			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, errors.New("Error response code " + string(resp.StatusCode))
			}

			return &ProxyQuality{
				ConnectTime: int(time.Now().Sub(startTime).Seconds()),
			}, nil

		},
		ValidateInterval:  600,
		maxValidateThread: 50,

		Filter: &AcceptAll{},

		Provider: provider.CreateMultiLoader(
			provider.CreateKuaidaili(),
			provider.CreateSslproxies(),
			provider.CreateIncloakk(),
		),

		ProxyStrategy: func(alive map[ProxyInfo]ProxyState, request *http.Request) *ProxyInfo {
			for info, _ := range alive {
				if request.URL.Scheme == "https" && info.Type&HTTPS == 0 {
					continue
				}
				return &info
			}
			return nil
		},
	}
}

type ValidateFunc func(ProxyInfo) (*ProxyQuality, error)

// TODO modifiable in outer package
type ProxyStrategy func(map[ProxyInfo]ProxyState, *http.Request) *ProxyInfo
