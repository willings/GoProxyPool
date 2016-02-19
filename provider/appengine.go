package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	client := p.client
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest("GET", "http://"+p.appid+".appspot.com/proxy.json", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Http Status code: " + strconv.Itoa(resp.StatusCode))
	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(resp.Body)

	ret := make([]*ProxyItem, 0)
	err = json.Unmarshal(buf.Bytes(), &ret)

	return ret, err
}
