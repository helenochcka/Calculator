package concurrent_map

import (
	"sync"
)

type ResultStorage struct {
	mu sync.Mutex
	m  map[string]int
}

func NewResultStorage() *ResultStorage {
	return &ResultStorage{m: make(map[string]int)}
}

func (rs *ResultStorage) Get(key string) *int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if value, ok := rs.m[key]; ok {
		return &value
	}
	return nil
}

func (rs *ResultStorage) Insert(key string, value int) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.m[key] = value
}
