package proxypool

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestReloadClose(t *testing.T) {
	var FILE_NAME = strconv.Itoa(int(time.Now().Unix()/1000)) + "test.csv"
	f, err := os.Create(FILE_NAME)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	defer os.Remove(FILE_NAME)

	pool, err := NewWithCsv(FILE_NAME, DefaultConfig())
	if err != nil {
		t.Error(err)
	}

	pool.Add(ProxyInfo{
		Host: "127.0.0.1",
		Port: 8080,
		Type: Protocol(1),
	})

	pool.Save()
}
