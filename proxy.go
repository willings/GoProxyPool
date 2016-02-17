package proxypool

import (
	"errors"
	"fmt"
	"github.com/willings/proxypool/storage"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Protocol int

const (
	HTTP  Protocol = 1
	HTTPS Protocol = 2

//SOCK5 Protocol = 4
)

type Proxy func(*http.Request) (*url.URL, error)

type ProxyInfo struct {
	Host      string
	Port      int
	Type      Protocol
	Anonymous bool
}

func (info *ProxyInfo) copyProxyInfo() *ProxyInfo {
	return &ProxyInfo{
		Host:      info.Host,
		Port:      info.Port,
		Type:      info.Type,
		Anonymous: info.Anonymous,
	}
}

func (info *ProxyInfo) GetProxy() Proxy {
	f := func(req *http.Request) (*url.URL, error) {
		if req.URL.Scheme == "https" && info.Type&HTTPS > 0 {
			return nil, errors.New("HTTPS is not supported")
		}
		return url.Parse(fmt.Sprintf("http://%s:%d", info.Host, info.Port))
	}
	return f
}

type ProxyState struct {
	InsertTime time.Time
	ActiveTime time.Time

	SuccessCnt int
	FailCnt    int

	Quality ProxyQuality
	AvgQ    ProxyQuality
}

func (state *ProxyState) copyProxyState() *ProxyState {
	return &ProxyState{
		InsertTime: state.InsertTime,
		ActiveTime: state.ActiveTime,
		SuccessCnt: state.SuccessCnt,
		FailCnt:    state.FailCnt,
		Quality:    *state.Quality.copyProxyQuality(),
		AvgQ:       *state.Quality.copyProxyQuality(),
	}
}

type ProxyQuality struct {
	ConnectTime   int
	DownloadSpeed int
}

func (q *ProxyQuality) copyProxyQuality() *ProxyQuality {
	return &ProxyQuality{
		ConnectTime:   q.ConnectTime,
		DownloadSpeed: q.DownloadSpeed,
	}
}

type ProxyPool struct {
	lock   sync.Locker
	active map[ProxyInfo]ProxyState

	config  *Config
	storage storage.ListStorage
}

func NewWithCsv(csvFilePath string, config *Config) (*ProxyPool, error) {
	if config == nil {
		config = DefaultConfig()
	}

	storage := storage.NewCsvStorage(csvFilePath)

	all, err := storage.Load()
	if err != nil {
		return nil, err
	}

	pool := &ProxyPool{
		lock:    &sync.Mutex{},
		active:  make(map[ProxyInfo]ProxyState),
		config:  config,
		storage: storage,
	}

	proxyInfos := make([]ProxyInfo, len(all))
	for _, entry := range all {
		info := ProxyInfo{
			Host:      entry.Host,
			Port:      entry.Port,
			Type:      Protocol(entry.Type),
			Anonymous: entry.Anonymous == 1,
		}
		if config.Filter.Accept(&info) {
			proxyInfos = append(proxyInfos, info)
		}
	}

	qualities := validateParallel(proxyInfos, config.Validator, 100)
	alive := make(map[ProxyInfo]ProxyState)
	for i, entry := range all {
		if qualities[i] != nil {
			alive[proxyInfos[i]] = ProxyState{
				InsertTime: entry.InsertTime,
				ActiveTime: time.Now(),
				SuccessCnt: entry.SuccessCnt,
				FailCnt:    entry.FailCnt,
				Quality:    *qualities[i],
			}
		}
	}

	pool.lock.Lock()
	pool.active = alive
	pool.lock.Unlock()

	return pool, nil
}

func (pool *ProxyPool) Save() {
	if pool.storage != nil {
		list := make([]*storage.ProxyEntry, len(pool.active))
		i := 0
		for info, state := range pool.active {
			list[i] = &storage.ProxyEntry{
				Host:      info.Host,
				Port:      info.Port,
				Type:      int(info.Type),
				Anonymous: 0,

				InsertTime: state.InsertTime,
				ActiveTime: state.ActiveTime,

				SuccessCnt: state.SuccessCnt,
				FailCnt:    state.FailCnt,

				LastConnectTime: state.Quality.ConnectTime,
				DownloadSpeed:   state.Quality.DownloadSpeed,
			}
			i++
		}

		pool.storage.Save(list)
	}
}

func (pool *ProxyPool) Reload() {
	if pool.config.Provider == nil {
		return
	}

	list, err := pool.config.Provider.Load()
	if err != nil {
		return
	}

	incoming := make([]ProxyInfo, len(list))
	for i, item := range list {
		incoming[i] = ProxyInfo{
			Host: item.Host,
			Port: item.Port,
			Type: Protocol(item.Type),
		}
	}

	qualities := validateParallel(incoming, pool.config.Validator, 100)

	active := make(map[ProxyInfo]ProxyState)
	for i, info := range incoming {
		if qualities[i] == nil {
			continue
		}
		active[info] = ProxyState{
			InsertTime: time.Now(),
			ActiveTime: time.Now(),

			SuccessCnt: 0,
			FailCnt:    0,

			Quality: *qualities[i],
		}
	}
	pool.lock.Lock()
	pool.active = active
	pool.lock.Unlock()
}

func (pool *ProxyPool) Count() int {
	return len(pool.active)
}

func (pool *ProxyPool) addAlive(pending []ProxyInfo, all []*storage.ProxyEntry) {

}

func (pool *ProxyPool) IsAvailable() bool {
	return len(pool.active) > 0
}

func (pool *ProxyPool) Add(info ProxyInfo) {
	postAddFunc := func() {
		if quality, err := pool.config.Validator(info); err == nil && quality != nil {
			pool.lock.Lock()
			pool.active[info] = ProxyState{
				InsertTime: time.Now(),
				ActiveTime: time.Now(),
				SuccessCnt: 0,
				FailCnt:    0,
				Quality:    *quality,
				AvgQ:       *quality,
			}
			pool.lock.Unlock()
		}
	}

	postAddFunc()
}

func (pool *ProxyPool) AutoProxy() Proxy {
	return func(req *http.Request) (*url.URL, error) {
		pool.lock.Lock()
		entry := pool.config.ProxyStrategy(pool.active, req)
		pool.lock.Unlock()

		if entry != nil {
			if req.URL.Scheme == "https" {
				if entry.Type&HTTPS > 0 {
					urlStr := fmt.Sprintf("https://%s:%d", entry.Host, entry.Port)
					return url.Parse(urlStr)
				} else {
					panic("Wrong stategy return: " + entry.Host + ":" + string(entry.Port))
				}
			} else {
				urlStr := fmt.Sprintf("http://%s:%d", entry.Host, entry.Port)
				return url.Parse(urlStr)
			}
		} else {
			if pool.config.AllowDirect {
				return nil, nil
			} else {
				return nil, errors.New("No Proxy available")
			}
		}
	}
}

func (pool *ProxyPool) GetProxy(requestURL *url.URL) *ProxyInfo {
	pool.lock.Lock()
	entry := pool.config.ProxyStrategy(pool.active, &http.Request{URL: requestURL})
	pool.lock.Unlock()
	return entry
}

func (pool *ProxyPool) GetProxyState(info ProxyInfo) *ProxyState {
	state := pool.active[info]
	return state.copyProxyState()
}

func (pool *ProxyPool) GetAllFixedProxy() []*ProxyInfo {
	ret := make([]*ProxyInfo, len(pool.active))
	pool.lock.Lock()
	i := 0
	for info, _ := range pool.active {
		ret[i] = info.copyProxyInfo()
		i++
	}
	pool.lock.Unlock()
	return ret
}

func (pool *ProxyPool) Feedback(proxy ProxyInfo, success bool, q *ProxyQuality) {
	pool.lock.Lock()

	state := pool.active[proxy]
	if success {
		state.Quality = *q
		state.AvgQ.DownloadSpeed = (state.SuccessCnt*state.AvgQ.DownloadSpeed + q.DownloadSpeed) / (state.SuccessCnt + 1)
		state.AvgQ.ConnectTime = (state.SuccessCnt*state.AvgQ.ConnectTime + q.ConnectTime) / (state.SuccessCnt + 1)
		state.SuccessCnt++
	} else {
		state.FailCnt++
	}

	pool.active[proxy] = state

	pool.lock.Unlock()
}
