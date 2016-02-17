package storage

import "time"

type ProxyEntry struct {
	Id        int
	Host      string
	Port      int
	Type      int
	Anonymous int

	InsertTime time.Time
	ActiveTime time.Time

	SuccessCnt int
	FailCnt    int

	LastConnectTime int
	DownloadSpeed   int
}

type ListStorage interface {
	Load() ([]*ProxyEntry, error)

	Save([]*ProxyEntry) error
}
