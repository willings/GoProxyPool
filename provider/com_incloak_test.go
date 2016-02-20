package provider

import (
	"fmt"
	"testing"
)

func TestIncloakLoad(t *testing.T) {
	p := CreateIncloakk()
	items, err := p.Load()
	fmt.Println("Loaded", len(items), "error:", err)
	if len(items) == 0 || err != nil {
		t.Fatal()
	}
}
