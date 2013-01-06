package main

import (
	"time"
	"encoding/json"
	"math/rand"
)

const idChars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const idSize = 3

var store Store = NewMemoryStore()

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
