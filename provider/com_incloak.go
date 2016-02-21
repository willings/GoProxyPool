package provider

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	IN_CLOAK_URL = "http://incloak.com/proxy-list"
	IN_CLOAK_PARAM = "type=hs&anon=234"
	IN_CLOAK_DEFAULT_PORTS = "80|8080|3128"
)

type Com_Incloak struct {
	Ports  []int
	client *http.Client
}

func (p *Com_Incloak) SetClient(client *http.Client) {
	p.client = client
}

func (p *Com_Incloak) Load() ([]*ProxyItem, error) {
	if p.Ports == nil || len(p.Ports) == 0{
		ports := strings.Split(IN_CLOAK_DEFAULT_PORTS, "|")
		p.Ports = make([]int, len(ports))
		for i, portStr := range ports {
			p.Ports[i], _ = strconv.Atoi(portStr)
		}

	}
	params := make([]interface{}, 0, len(p.Ports))
	for _, port := range p.Ports {
		params = append(params, port)
	}

	return loadParallel(p, 5, params...)

}

func (p *Com_Incloak) load(param interface{}) ([]*ProxyItem, error) {
	port, found := param.(int)
	if !found {
		return nil, errors.New("Wrong param")
	}

	url := IN_CLOAK_URL + "/?ports=" + strconv.Itoa(port) + "&" + IN_CLOAK_PARAM

	client := p.client
	if client == nil {
		client = http.DefaultClient
	}

	buf, err := httpGet(url, p.client)
	if err != nil {
		return nil, err
	}

	startBytes := []byte("<table class=pl ")
	endBytes := []byte("</table>")

	tblStart := bytes.Index(buf, startBytes)
	if tblStart < 0 {
		return nil, errors.New("Failed to parse stream")
	}
	tblEnd := bytes.Index(buf[tblStart:], endBytes) + tblStart
	if tblEnd <= tblStart {
		return nil, errors.New("Failed to parse stream")
	}

	b := buf[tblStart : tblEnd+len(endBytes)]
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
