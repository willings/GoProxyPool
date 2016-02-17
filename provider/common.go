package provider

type loadUnit interface {
	load(param interface{}) ([]*ProxyItem, error)
}

type result struct {
	items []*ProxyItem
	err   error
}

func loadParallel(unit loadUnit, maxRoutine int, params ...interface{}) ([]*ProxyItem, error) {
	results := make(chan *result)
	routines := make(chan int, maxRoutine)

	go func() {
		for i, param := range params {
			routines <- i
			go func(p interface{}) {
				items, err := unit.load(p)
				results <- &result{
					items: items,
					err:   err,
				}
				<-routines
			}(param)
		}
	}()

	ret := make([]*ProxyItem, 0)
	var lastError error

	for i := 0; i < len(params); i++ {
		result := <-results
		if result.items != nil {
			ret = append(ret, result.items...)
		}
		if result.err != nil {
			lastError = result.err
		}
	}

	return ret, lastError
}

type proxyKey struct {
	Host string
	Port int
}

type proxyMeta struct {
	Type      int
	Anonymous int
}

func removeDuplicates(items []*ProxyItem) []*ProxyItem {
	m := make(map[proxyKey]proxyMeta)
	for _, item := range items {
		key := proxyKey{
			Host: item.Host,
			Port: item.Port,
		}
		if v, ok := m[key]; ok {
			m[key] = proxyMeta{
				Type:      item.Type,
				Anonymous: item.Anonymous,
			}
		} else {
			m[key] = proxyMeta{
				Type:      item.Type | v.Type,
				Anonymous: item.Anonymous & v.Anonymous,
			}
		}
	}

	ret := make([]*ProxyItem, 0, len(m))
	for k, v := range m {
		ret = append(ret, &ProxyItem{
			Host:      k.Host,
			Port:      k.Port,
			Type:      v.Type,
			Anonymous: v.Anonymous,
		})
	}
	return ret
}
