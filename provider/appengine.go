package provider

import (
	"encoding/json"
	"net/http"
)

type AppEngine struct {
	appid  string
	client *http.Client
}

func CreateAppEngine() *AppEngine {
	return &AppEngine{
		appid: "free-proxy-list",
	}
}

func (p *AppEngine) Load() ([]*ProxyItem, error) {
	return p.load(nil)
}

func (p *AppEngine) SetClient(client *http.Client) {
	p.client = client
}

func (p *AppEngine) load(config interface{}) ([]*ProxyItem, error) {
	b, err := httpGet("http://"+p.appid+".appspot.com/refresh", p.client)
	if err != nil {
		return nil, err
	}

	ret := make([]*ProxyItem, 0)
	err = json.Unmarshal(b, &ret)

	return ret, err
}
