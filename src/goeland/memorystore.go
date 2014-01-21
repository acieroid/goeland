package main

import (
	"sync"
)

type MemoryStore struct {
	lists map[string] *TodoList
	mu sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		make(map[string] *TodoList),
		sync.RWMutex{},
	}
}

func (s *MemoryStore) Get(id string) *TodoList {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lists[id]
}

func (s *MemoryStore) Exists(id string) bool {
	_, exists := s.lists[id]
	return exists
}

func (s *MemoryStore) Set(list *TodoList) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lists[list.Id] = list
	return true
}
