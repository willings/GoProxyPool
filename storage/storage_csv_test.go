package storage

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	const FILE = "test.csv.tmp"

	storage := NewCsvStorage(FILE)
	if storage == nil {
		t.Fail()
	}

	all, err := storage.Load()
	if err != nil || len(all) != 186 {
		t.Fatal()
	}
}

func TestSaveLoad(t *testing.T) {
	const FILE = "test.csv"
	const N = 10

	storage := NewCsvStorage(FILE)
	if storage == nil {
		t.Fail()
	}

	list := make([]*ProxyEntry, N)
	for i := 0; i < N; i++ {
		list[i] = randEntry()
	}
	storage.Save(list)

	out, err := storage.Load()

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(list, out) {
		t.Fail()
	}

	os.Remove(FILE)
}

func randEntry() *ProxyEntry {
	return &ProxyEntry{
		Id:        rand.Intn(100),
		Host:      randHost(),
		Port:      rand.Intn(65535),
		Type:      rand.Intn(4),
		Anonymous: rand.Intn(2),

		InsertTime: time.Unix(time.Now().Unix(), 0),
		ActiveTime: time.Unix(time.Now().Unix(), 0),

		SuccessCnt: rand.Intn(65535),
		FailCnt:    rand.Intn(65535),

		LastConnectTime: rand.Intn(10000),
		DownloadSpeed:   rand.Intn(10000),
	}
}

func randHost() string {
	i1 := rand.Intn(255)
	i2 := rand.Intn(255)
	i3 := rand.Intn(255)
	i4 := rand.Intn(255)

	return fmt.Sprintf("%d.%d.%d.%d", i1, i2, i3, i4)

}
