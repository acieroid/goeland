package main

import (
	"time"
	"sync"
	"encoding/json"
	"math/rand"
)

const idChars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const idSize = 3

type TodoListItem struct {
	Name string
	Description string
	Status string
}

type TodoList struct {
	Id string
	Name string
	ModificationTime int64
	Items[] *TodoListItem
}

type TodoListStore struct {
	lists map[string] *TodoList
	mu sync.RWMutex
}

var store *TodoListStore = NewTodoListStore()

func NewTodoListStore() *TodoListStore {
	return &TodoListStore{
		make(map[string] *TodoList),
		sync.RWMutex{},
	}
}

func (s *TodoListStore) Get(id string) *TodoList {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lists[id]
}

func (s *TodoListStore) Exists(id string) bool {
	_, exists := s.lists[id]
	return exists
}

func (s *TodoListStore) Set(list *TodoList) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lists[list.Id] = list
	return true
}

func RandomId() string {
	buf := make([]byte, idSize)
	for i := 0; i < idSize; i++ {
		buf[i] = idChars[rand.Intn(len(idChars)-1)]
	}
	return string(buf)
}

func Now() int64 {
	return time.Now().UTC().UnixNano()
}

func NewList(name string) *TodoList {
	var id string
	for id = RandomId(); ListExists(id); id = RandomId() {
	}
	return &TodoList{id, name, Now(), nil}
}

func LoadList(id string) *TodoList {
	return store.Get(id)
}

func ParseList(descr []byte) *TodoList {
	l := &TodoList{}
	json.Unmarshal(descr, l)
	return l
}

func ListExists(id string) bool {
	return store.Exists(id)
}

func (l *TodoList) Save() bool {
	return store.Set(l)
}

func (l *TodoList) AddItem(item *TodoListItem) {
	l.Items = append(l.Items, item)
}

func InitModel() {
	rand.Seed(Now());
}
