package provider

import (
	"encoding/json"
	"net/http"
)

type Com_appspot struct {
	Appid  string
	client *http.Client
}

func (p *Com_appspot) Load() ([]*ProxyItem, error) {
	return p.load(nil)
}

func (p *Com_appspot) SetClient(client *http.Client) {
	p.client = client
}

func (p *Com_appspot) load(config interface{}) ([]*ProxyItem, error) {
	b, err := httpGet("http://"+p.Appid +".appspot.com/refresh", p.client)
	if err != nil {
		return nil, err
	}

	ret := make([]*ProxyItem, 0)
	err = json.Unmarshal(b, &ret)

	return ret, err
}
