package main

import (
	"time"
	"encoding/json"
	"math/rand"
)

const idChars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const idSize = 3

var lists = map[string] *TodoList {};

type Message struct {
	Message string
	Action string
	ActionURL string
}

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

func randomId() string {
	buf := make([]byte, idSize)
	for i := 0; i < idSize; i++ {
		buf[i] = idChars[rand.Intn(len(idChars)-1)]
	}
	return string(buf)
}

func now() int64 {
	return time.Now().UTC().UnixNano()
}

func newList(name string) *TodoList {
	var id string
	for id = randomId(); listExists(id); id = randomId() {
	}
	return &TodoList{id, name, now(), nil}
}

func loadList(id string) *TodoList {
	return lists[id]
}

func parseList(descr []byte) *TodoList {
	l := &TodoList{}
	json.Unmarshal(descr, l)
	return l
}

func listExists(name string) bool {
	_, exists := lists[name]
	return exists
}

func init() {
	rand.Seed(now());
}

func (l *TodoList) save() {
	lists[l.Id] = l
}

func (l *TodoList) addItem(item *TodoListItem) {
	l.Items = append(l.Items, item)
}
