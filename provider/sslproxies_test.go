package provider

import (
	"fmt"
	"testing"
)

func TestSslProxiesLoad(t *testing.T) {
	provider := CreateSslproxies()
	list, err := provider.Load()
	fmt.Println("loaded", len(list), "error", err)
	if err != nil {
		t.Fail()
	}
}
