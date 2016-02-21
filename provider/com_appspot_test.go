package provider

import (
	"testing"
)

func TestAppEngineLoad(t *testing.T) {
	provider := &Com_appspot{
		Appid: "free-proxy-list",
	}
	assertProxyExists(provider, t)
}
