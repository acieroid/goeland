package main

import (
	"code.google.com/p/gosqlite/sqlite"
	"log"
	"os"
	"runtime"
)

func CheckSQLError(err error) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("Error with SQLite: %s, at %s:%d\n", err, file, line)
		return false
	}
	return true
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

type SQLiteStore struct {
	conn *sqlite.Conn
}

func NewSQLiteStore(filename string) *SQLiteStore {
	s := &SQLiteStore{}
	log.Println("Opening SQLite store in file", filename)

	defer s.CreateTables()

	/* open the connection */
	conn, err := sqlite.Open(filename)
	if err != nil {
		log.Fatal("Cannot open sqlite connection:", err)
	}
	s.conn = conn

	return s
}

func (s *SQLiteStore) CreateTable(descr string) {
	err := s.conn.Exec(descr)
	if err != nil {
		log.Fatal("Error when creating SQLite tables:", err)
		return
	}
}

func (s *SQLiteStore) CreateTables() {
	s.CreateTable("create table if not exists list (" +
		"id integer primary key autoincrement," +
		"idstr text unique," +
		"name text," +
		"mtime int(64))")
	s.CreateTable("create table if not exists item (" +
		"list_id integer not null," +
		"name text," +
		"descr text," +
		"status text," +
		"foreign key (list_id) references list(id) on delete cascade)")
}

func (s *SQLiteStore) Get(idstr string) *TodoList {
	var id int
	l := &TodoList{}
	stmt, err := s.conn.Prepare("select id, idstr, name, mtime from list where idstr = ?")
	if !CheckSQLError(err) {
		return nil
	}

	err = stmt.Exec(idstr)
	if !CheckSQLError(err) {
		return nil
	}

	err = stmt.Scan(&id, &l.Id, &l.Name, &l.ModificationTime)
	if !CheckSQLError(err) {
		return nil
	}

	stmt, err = s.conn.Prepare("select name, descr, status from item where list_id = ?")
	if !CheckSQLError(err) {
		return nil
	}

	err = stmt.Exec(id)
	if !CheckSQLError(err) {
		return nil
	}

	for stmt.Next() {
		item := &TodoListItem{}
		err = stmt.Scan(&item.Name, &item.Description, &item.Status)
		if !CheckSQLError(err) {
			return nil
		}
		l.AddItem(item)
	}
	return l
}

func (s *SQLiteStore) Exists(id string) bool {
	var n int
	stmt, err := s.conn.Prepare("select count(*) from list where idstr = ?")
	if !CheckSQLError(err) {
		return false
	}

	err = stmt.Exec(id)
	if !CheckSQLError(err) {
		return false
	}

	err = stmt.Scan(&n)
	if !CheckSQLError(err) {
		return false
	}

	return n != 0
}

func (s *SQLiteStore) Set(list *TodoList) bool {
	var id int
	err := s.conn.Exec("insert into list(idstr, name, mtime) values (?, ?, ?)", list.Id, list.Name, list.ModificationTime)
	if !CheckSQLError(err) {
		return false
	}

	stmt, err := s.conn.Prepare("select id from list where idstr = ?")
	if !CheckSQLError(err) {
		return false
	}

	err = stmt.Exec(list.Id)
	if !CheckSQLError(err) {
		return false
	}

	err = stmt.Scan(&id)
	if !CheckSQLError(err) {
		return false
	}

	for _, item := range list.Items {
		stmt, err := s.conn.Prepare("insert into item(list_id, name, descr, status) values (?, ?, ?, ?)")
		if !CheckSQLError(err) {
			return false
		}

		err = stmt.Exec(id, item.Name, item.Description, item.Status)
		if !CheckSQLError(err) {
			return false
		}
	}

	return true
}
