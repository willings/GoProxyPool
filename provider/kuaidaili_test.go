package provider

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	provider := CreateKuaidaili()
	list, err := provider.Load()
	fmt.Println("loaded", len(list), "error", err)
	if err != nil && len(list) == 0 {
		t.Fail()
	}
}
