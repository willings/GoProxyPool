package provider

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func httpGet(url string, client *http.Client) ([]byte, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client.Timeout = time.Second * 10

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = errors.New("Response code is " + strconv.Itoa(resp.StatusCode))
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)

	return buf.Bytes(), err
}

func addBotHeader(h http.Header) {
	h.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	h.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	h.Add("Accept-Encoding", "gzip, deflate, sdch")
	h.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4")
	h.Add("Cache-control", "upgrade-insecure-requests")
	h.Add("Upgrade-Insecure-Requests", "1")
	h.Add("Connection", "Keep-Alive")
}
