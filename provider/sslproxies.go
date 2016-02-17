package provider

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Org_Sslproxies struct {
	url    string
	client *http.Client
}

func CreateSslproxies() *Org_Sslproxies {
	return &Org_Sslproxies{
		url: "http://www.sslproxies.org/",
	}
}

func (p *Org_Sslproxies) SetClient(client *http.Client) {
	p.client = client
}

func (p *Org_Sslproxies) Load() ([]*ProxyItem, error) {
	req, err := http.NewRequest("GET", p.url, nil)
	if err != nil {
		return nil, err
	}

	addBotHeader(req.Header)

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

	return ret, nil
}

func (p *Org_Sslproxies) convert(tr *Tr) *ProxyItem {
	if len(tr.Td) < 4 {
		return nil
	}

	port, err := strconv.Atoi(tr.Td[1])
	if err != nil || port == 0 {
		return nil
	}

	var t int
	if strings.Contains(tr.Td[6], "HTTPS") {
		t = 3
	} else if strings.Contains(tr.Td[6], "HTTP") {
		t = 1
	}

	var a int
	if tr.Td[4] == "anonymous" {
		a = 1
	} else {
		a = 0
	}

	return &ProxyItem{
		Host:      tr.Td[0],
		Port:      port,
		Type:      t,
		Anonymous: a,
	}
}
