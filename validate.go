package proxypool

type verifyResult struct {
	i       int
	info    *ProxyInfo
	quality *ProxyQuality
}

func validateParallel(pending []ProxyInfo, vFunc ValidateFunc, limit int) []*ProxyQuality {
	c := make(chan *verifyResult)
	f := func(i int, info *ProxyInfo) {
		quality, _ := vFunc(*info)
		c <- &verifyResult{
			i:       i,
			info:    info,
			quality: quality,
		}
	}

	sem := make(chan int, limit)
	for i, _ := range pending {
		sem <- 1
		go f(i, &pending[i])
		<-sem
	}

	qualities := make([]*ProxyQuality, len(pending))
	for i := 0; i < len(pending); i++ {
		result := <-c
		qualities[result.i] = result.quality
	}
	return qualities
}
