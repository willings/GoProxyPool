package proxypool

type Filter interface {
	Accept(info *ProxyInfo) bool
}

type AcceptAll struct {
}

func (filter *AcceptAll) Accept(info *ProxyInfo) bool {
	return true
}

type AcceptOnlyAnnoymous struct {
}

func (filter *AcceptOnlyAnnoymous) Accept(info *ProxyInfo) bool {
	return info.Anonymous
}
