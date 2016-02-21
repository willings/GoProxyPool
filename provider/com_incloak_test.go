package provider

import (
	"testing"
)

func TestIncloakLoad(t *testing.T) {
	p := &Com_Incloak{}
	assertProxyExists(p, t)
}
