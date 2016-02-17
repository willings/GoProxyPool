package storage

type MemStorage struct {
	EntryList []*ProxyEntry
}

func (storage *MemStorage) Load() ([]*ProxyEntry, error) {
	return storage.EntryList, nil
}

func (storage *MemStorage) Save(entryList []*ProxyEntry) error {
	storage.EntryList = entryList
	return nil
}
