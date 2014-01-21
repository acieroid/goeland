package main

type Store interface {
	Get(id string) *TodoList
	Exists(id string) bool
	Set(list *TodoList) bool
}