package provider

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
)

type Com_Incloak struct {
	url    string
	param  string
	ports  []int
	client *http.Client
}

func CreateIncloakk() *Com_Incloak {
	ports := make([]int, 0)
	ports = append(ports, 80, 8080, 3128)
	return &Com_Incloak{
		url:   "http://incloak.com/proxy-list",
		param: "type=hs&anon=234",
		ports: ports,
	}
}

func (p *Com_Incloak) SetClient(client *http.Client) {
	p.client = client
}

func (p *Com_Incloak) Load() ([]*ProxyItem, error) {
	params := make([]interface{}, 0, len(p.ports))
	for _, port := range p.ports {
		params = append(params, port)
	}

	return loadParallel(p, 5, params...)

}

func (p *Com_Incloak) load(param interface{}) ([]*ProxyItem, error) {
	port, found := param.(int)
	if !found {
		return nil, errors.New("Wrong param")
	}

	url := p.url + "/?ports=" + strconv.Itoa(port) + "&" + p.param

	client := p.client
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// addBotHeader(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Http Status code: " + strconv.Itoa(resp.StatusCode))
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read stream")
	}

	startBytes := []byte("<table class=pl ")
	endBytes := []byte("</table>")

	tblStart := bytes.Index(buf.Bytes(), startBytes)
	if tblStart < 0 {
		return nil, errors.New("Failed to parse stream")
	}
	tblEnd := bytes.Index(buf.Bytes()[tblStart:], endBytes) + tblStart
	if tblEnd <= tblStart {
		return nil, errors.New("Failed to parse stream")
	}

	b := buf.Bytes()[tblStart : tblEnd+len(endBytes)]
	trArr := bytes.Split(b, []byte("<tr>"))
	ret := make([]*ProxyItem, len(trArr))
	cnt := 0
	for _, tr := range trArr {
		item := p.convert(tr, port)
		if item != nil {
			ret[cnt] = item
			cnt++
		}
	}
	if cnt > 0 {
		return ret[0:cnt], nil
	} else {
		return nil, errors.New("Could not find proxy")
	}

	return nil, nil
}

func (p *Com_Incloak) convert(tr []byte, port int) *ProxyItem {
	tdBytes := []byte("<td class=tdl>")
	if bytes.Index(tr, tdBytes) != 0 {
		return nil
	}
	tdCloseBytes := []byte("</td>")
	tdClosePos := bytes.Index(tr, tdCloseBytes)
	host := string(tr[len(tdBytes):tdClosePos])

	var t int
	if bytes.Index(tr, []byte("HTTPS")) > 0 {
		t = 3
	} else {
		t = 1
	}

	return &ProxyItem{
		Host:      host,
		Port:      port,
		Type:      t,
		Anonymous: 1,
	}
}
