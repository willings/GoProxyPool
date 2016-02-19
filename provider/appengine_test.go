package provider

import (
	"fmt"
	"testing"
)

func TestAppEngineLoad(t *testing.T) {

	provider := CreateAppEngine()
	items, err := provider.Load()
	fmt.Println("Loaded", len(items), "error:", err)
	if len(items) == 0 || err != nil {
		t.Fatal()
	}
}
