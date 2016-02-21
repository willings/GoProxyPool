package provider

import "net/http"

type MultiLoader struct {
	l []ProxyProvider
}

func CreateProvider(providerName string) ProxyProvider {
	switch providerName {
	case "Kuaidaili":
		return &Com_kuaidaili{}

	case "Sslproxies":
		return &Org_sslproxies{}

	case "Incloak":
		return &Com_Incloak{
			Ports: []int{80, 8080, 3128},
		}

	default:
		return nil
	}
}

func CreateAllProvider() []ProxyProvider {
	return []ProxyProvider{
		&Com_kuaidaili{
			Page: 10,
		},
		&Org_sslproxies{},
		&Com_Incloak{Ports: []int{80, 8080, 3128},},
	}
}

func CreateAllLoader() ProxyProvider {
	ret := &MultiLoader{
		l: CreateAllProvider(),
	}
	return ret
}

func CreateMultiLoader(l ...ProxyProvider) ProxyProvider {
	ret := &MultiLoader{
		l: l,
	}

	return ret
}

func (loader *MultiLoader) Load() ([]*ProxyItem, error) {
	var lastError error
	ret := make([]*ProxyItem, 0)
	for _, ld := range loader.l {
		items, err := ld.Load()
		if err != nil {
			lastError = err
		}
		ret = append(ret, items...)
	}
	if len(ret) == 0 {
		return ret, lastError
	} else {
		return removeDuplicates(ret), nil
	}
}

func (loader *MultiLoader) SetClient(client *http.Client) {
	for _, ld := range loader.l {
		ld.SetClient(client)
	}
}
