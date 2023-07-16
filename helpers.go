package goentdb

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
)

func MD5(i string) string {
	data := []byte(i)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func LoadMapFromFilepath(filepath string, dict *map[int]*EntKeyword, lock sync.RWMutex) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	lock.Lock()
	defer lock.Unlock()

	decoder := gob.NewDecoder(f)
	decoder.Decode(dict)

	return nil
}

func DumpMapToFilepath(filepath string, dict map[int]*EntKeyword, lock sync.RWMutex) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(f)

	lock.RLock()
	defer lock.RUnlock()

	if err := encoder.Encode(dict); err != nil {
		return err
	}
	f.Close()
	return nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
