package provider

import (
	"testing"
)

func TestLoad(t *testing.T) {
	provider := &Com_kuaidaili{
		Page: 10,
	}
	assertProxyExists(provider, t)
}
