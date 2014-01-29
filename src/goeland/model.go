package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

const idChars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const idSize = 3

var store Store = NewSQLiteStore("goeland.db")

// XXX add root hashtable to handle every TodoList

type TodoList struct {
	Id               int64
	Name             string
	Description      string
	Status           string
	ModificationTime int64
	EstimatedTime    int64
	Items            []int64
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
	return &TodoList{id, name, "", "none", Now(), 0, nil}
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

func SumEtime(items[] int64) {
	if items == nil {
		return l.EstimatedTime;
	}
	for _, item := range l.Items {
		if l.items != nil
			return l.EstimatedTime;
		l.EstimatedTime = SumEtime(l.items)
		l.EstimatedTime += item.EstimatedTime
	}
}

func (l *TodoList) Save() bool {
	l.SumEtime()
	return store.Set(l)
}

func (l *TodoList) AddItem(item *TodoListItem) {
	l.Items = append(l.Items, item)
}

func FindItem(root[] *TodoList, id int64) *TodoList {
	for _, item := range items {
		if item.Id == id {
			return item
		}
		subItem := item.GetItem(id)
		if subItem != nil {
			return subItem
		}
	}
	return nil
}

func (l *TodoList) GetItem(id int64) *TodoListItem {
	return FindItem(l.Items, id)
}

func (i *TodoListItem) GetItem(id int64) *TodoListItem {
	return FindItem(i.Items, id)
}

func (i *TodoListItem) AddItem(item *TodoListItem) {
	i.Items = append(i.Items, item)
}		

func InitModel() {
	rand.Seed(Now())
}
