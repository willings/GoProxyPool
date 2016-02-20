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
	addBotHeader(req.Header)

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
	h.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.54 Safari/537.36")
	h.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	// h.Add("Accept-Encoding", "gzip, deflate, sdch")
	h.Add("Accept-Language", "en-US;q=0.6,en;q=0.4")
	h.Add("Cache-control", "no-cache")
	h.Add("Upgrade-Insecure-Requests", "1")
	h.Add("Connection", "Keep-Alive")
	h.Add("Cookie", `visited=2016%2F02%2F21+00%3A37%3A03; __utmt=1; hl=en; pv=5; userno=20160221-000034; from=direct; __atuvc=4%7C8; __atuvs=56c8881f3c6961e7003; __utma=127656268.820237221.1455982624.1455982624.1455982624.1; __utmb=127656268.8.10.1455982624; __utmc=127656268; __utmz=127656268.1455982624.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmv=127656268.Japan`)
}
