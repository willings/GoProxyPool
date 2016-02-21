package provider

import (
	"testing"
	"fmt"
)

func assertProxyExists(provider ProxyProvider, t *testing.T) {

	items, err := provider.Load()
	fmt.Println("Loaded", len(items), "error:", err)
	if len(items) == 0 || err != nil {
		t.Fatal()
	}

}
