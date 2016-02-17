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

type Com_Kuaidaili struct {
	urlbase string
	query   []string
	page    int
	client  *http.Client
}

func CreateKuaidaili() *Com_Kuaidaili {
	return &Com_Kuaidaili{
		urlbase: "http://www.kuaidaili.com/free/",
		query:   []string{"inha", "intr", "outha", "outtr"},
		page:    10,
	}
}

func (p *Com_Kuaidaili) SetClient(client *http.Client) {
	p.client = client
}

func (p *Com_Kuaidaili) Load() ([]*ProxyItem, error) {
	N := len(p.query) * p.page
	params := make([]interface{}, 0, N)
	for _, q := range p.query {
		for i := 1; i <= p.page; i++ {
			url := p.urlbase + q + "/" + strconv.Itoa(i)
			params = append(params, url)
		}
	}

	return loadParallel(p, 10, params...)
}

func (p *Com_Kuaidaili) load(param interface{}) ([]*ProxyItem, error) {
	url, found := param.(string)
	if !found {
		return nil, errors.New("Wrong params type")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// addBotHeader(req.Header)

	client := p.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Http Status code: " + string(resp.StatusCode))
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read stream")
	}

	startBytes := []byte("<tbody>")
	endBytes := []byte("</tbody>")

	tbodyStart := bytes.Index(buf.Bytes(), startBytes)
	tbodyEnd := bytes.Index(buf.Bytes(), endBytes)
	if tbodyEnd <= tbodyStart {
		return nil, errors.New("Failed to parse stream")
	}

	bytes := buf.Bytes()[tbodyStart : tbodyEnd+len(endBytes)]
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

func (p *Com_Kuaidaili) convert(tr *Tr) *ProxyItem {
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
