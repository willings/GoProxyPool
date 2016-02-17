package proxypool

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Provider == nil {
		t.Fail()
	}
	if config.ProxyStrategy == nil {
		t.Fail()
	}
	if config.Validator == nil {
		t.Fail()
	}
}
