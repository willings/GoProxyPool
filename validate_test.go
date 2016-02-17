package proxypool

import (
	"testing"
	"time"
)

func TestValidateParallel(t *testing.T) {
	cfg := DefaultConfig()
	pList, err := cfg.Provider.Load()
	if err != nil {
		t.Fail()
	}
	if len(pList) == 0 {
		return
	}

	pending := make([]ProxyInfo, len(pList))
	for i, item := range pList {
		pending[i] = ProxyInfo{
			Host: item.Host,
			Port: item.Port,
			Type: Protocol(item.Type),
		}
	}

	c := make(chan []*ProxyQuality)
	go func() {
		c <- validateParallel(pending, cfg.Validator, 5)
	}()

	select {
	case qualities := <-c:
		if len(qualities) != len(pending) {
			t.Error("Missing varify result")
		}
	case <-time.After(time.Second * 12):
		t.Error("Timeout in validate")
	}

}
