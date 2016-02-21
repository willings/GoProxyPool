package provider

import "testing"

func TestCreateAllLoader(t *testing.T) {
	provider := CreateAllLoader()
	assertProxyExists(provider, t)
}
