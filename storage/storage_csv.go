package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CsvStorage struct {
	filePath string
}

func NewCsvStorage(filePath string) *CsvStorage {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, _ := os.Create(filePath)
		if file != nil {
			file.Close()
		}
	}

	return &CsvStorage{
		filePath: filePath,
	}
}

func (storage *CsvStorage) Load() ([]*ProxyEntry, error) {
	file, err := os.Open(storage.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]*ProxyEntry, 0)
	bufReader := bufio.NewReader(file)
	for {
		line, _, err := bufReader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}
		entry := parseLine(string(line))
		if entry != nil {
			ret = append(ret, entry)
		}
	}

	return ret, nil
}

func (storage *CsvStorage) Save(entries []*ProxyEntry) error {
	file, err := os.OpenFile(storage.filePath, os.O_RDWR, 0664)
	if err != nil {
		return nil
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	for i, entry := range entries {
		if entry.Id <= 0 {
			entry.Id = i
		}

		line := fmt.Sprintf(`%d,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d`,
			entry.Id, entry.Host, entry.Port, entry.Type, entry.Anonymous,
			entry.InsertTime.Unix(), entry.ActiveTime.Unix(),
			entry.SuccessCnt, entry.FailCnt,
			entry.LastConnectTime, entry.DownloadSpeed)
		bufWriter.WriteString(line)
		bufWriter.WriteString("\n")
	}
	bufWriter.Flush()
	return nil
}

func parseLine(line string) *ProxyEntry {
	elements := strings.Split(line, ",")
	if len(elements) != 11 {
		return nil
	}

	return &ProxyEntry{
		Id:        atoi(elements[0]),
		Host:      elements[1],
		Port:      atoi(elements[2]),
		Type:      atoi(elements[3]),
		Anonymous: atoi(elements[4]),

		InsertTime: time.Unix(int64(atoi(elements[5])), 0),
		ActiveTime: time.Unix(int64(atoi(elements[6])), 0),

		SuccessCnt: atoi(elements[7]),
		FailCnt:    atoi(elements[8]),

		LastConnectTime: atoi(elements[9]),
		DownloadSpeed:   atoi(elements[10]),
	}
}

func atoi(s string) int {
	ret, err := strconv.Atoi(s)
	if err == nil {
		return ret
	} else {
		return 0
	}
}
