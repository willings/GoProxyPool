package provider

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	KUAIDAILI_URL = "http://www.kuaidaili.com/free/"
	KUAIDAILI_PARAM = "inha|intr|outha|outtr"
	KUAIDAILI_PAGE = 10
)

type Com_kuaidaili struct {
	Page int
	client  *http.Client
}

func (p *Com_kuaidaili) SetClient(client *http.Client) {
	p.client = client
}

func (p *Com_kuaidaili) Load() ([]*ProxyItem, error) {
	if p.Page <= 0 {
		p.Page = KUAIDAILI_PAGE
	}
	queries := strings.Split(KUAIDAILI_PARAM, "|")
	N := len(queries) * p.Page
	params := make([]interface{}, 0, N)
	for _, q := range queries {
		for i := 1; i <= p.Page; i++ {
			url := KUAIDAILI_URL + q + "/" + strconv.Itoa(i)
			params = append(params, url)
		}
	}

	return loadParallel(p, 10, params...)
}

func (p *Com_kuaidaili) load(param interface{}) ([]*ProxyItem, error) {
	url, found := param.(string)
	if !found {
		return nil, errors.New("Wrong params type")
	}

	b, err := httpGet(url, p.client)
	if err != nil {
		return nil, errors.New("Failed to read stream")
	}

	startBytes := []byte("<tbody>")
	endBytes := []byte("</tbody>")

	tbodyStart := bytes.Index(b, startBytes)
	tbodyEnd := bytes.Index(b, endBytes)
	if tbodyEnd <= tbodyStart {
		return nil, errors.New("Failed to parse stream")
	}

	bytes := b[tbodyStart : tbodyEnd+len(endBytes)]
	tbl := Tbody{}
	err = xml.Unmarshal(bytes, &tbl)
	if err != nil {
		return nil, err
	}
	ret := make([]*ProxyItem, len(tbl.Tr))
	cnt := 0
	for _, tr := range tbl.Tr {
		item := p.convert(&tr)
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

func (p *Com_kuaidaili) convert(tr *Tr) *ProxyItem {
	if len(tr.Td) < 4 {
		return nil
	}

	port, err := strconv.Atoi(tr.Td[1])
	if err != nil || port == 0 {
		return nil
	}

	var t int
	if strings.Contains(tr.Td[3], "HTTPS") {
		t = 3
	} else if strings.Contains(tr.Td[3], "HTTP") {
		t = 1
	}

	var a int
	if tr.Td[2] == "高匿名" {
		a = 1
	} else {
		a = 0
	}

	_, err = url.Parse(fmt.Sprintf("http://%s:%d", tr.Td[0], port))
	if err != nil {
		return nil
	}

	return &ProxyItem{
		Host:      tr.Td[0],
		Port:      port,
		Type:      t,
		Anonymous: a,
	}
}
