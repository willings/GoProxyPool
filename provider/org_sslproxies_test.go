package provider

import (
	"testing"
)

func TestSslProxiesLoad(t *testing.T) {
	provider := &Org_sslproxies{}
	assertProxyExists(provider, t)
}
