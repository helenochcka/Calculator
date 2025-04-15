package concurrent_map

import (
	"Calculator/core"
	"sync"
)

type InstructionStorage struct {
	mu sync.Mutex
	m  map[string][]core.Instruction
}

func NewInstructionStorage() *InstructionStorage {

	return &InstructionStorage{m: make(map[string][]core.Instruction)}
}

func (is *InstructionStorage) Get(key string) *[]core.Instruction {
	is.mu.Lock()
	defer is.mu.Unlock()
	if values, ok := is.m[key]; ok {
		return &values
	}
	return nil
}

func (is *InstructionStorage) Insert(key string, value core.Instruction) {
	is.mu.Lock()
	defer is.mu.Unlock()
	is.m[key] = append(is.m[key], value)
}

func (is *InstructionStorage) Delete(key string) {
	is.mu.Lock()
	defer is.mu.Unlock()
	delete(is.m, key)
}
