package provider

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	provider := CreateKuaidaili()
	list, err := provider.Load()
	fmt.Println("loaded", list, "error", err)
	if err != nil {
		t.Fail()
	}
}
