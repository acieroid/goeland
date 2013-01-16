package main

import (
	"github.com/kuroneko/gosqlite3"
	"log"
	"os"
	"runtime"
)

func CheckSQLError(err error) bool {
	if err != nil && err != sqlite3.ROW {
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
	db *sqlite3.Database
}

func NewSQLiteStore(filename string) *SQLiteStore {
	s := &SQLiteStore{}
	log.Println("Opening SQLite store in file", filename)

	defer s.CreateTables()

	// open the connection
	db, err := sqlite3.Open(filename)
	if err != nil {
		log.Fatal("Cannot open sqlite connection:", err)
	}
	s.db = db

	return s
}

func (s *SQLiteStore) CreateTable(descr string) {
	_, err := s.db.Execute(descr)
	if err != nil {
		log.Fatal("Error when creating SQLite tables:", err)
		return
	}
}

func (s *SQLiteStore) CreateTables() {
	s.CreateTable(`create table if not exists list (
		       id integer primary key autoincrement,
		       idstr text unique,
		       name text,
		       mtime int(64))`)
	s.CreateTable(`create table if not exists item (
		       id integer primary key autoincrement,
		       parent_id integer not null,
		       list_id integer not null,
		       name text,
		       descr text,
		       status text,
		       foreign key (list_id) references list(id) on delete cascade,
		       foreign key (parent_id) references item(id) on delete cascade)`)
}

func (s *SQLiteStore) Get(idstr string) *TodoList {
	l := &TodoList{}
	stmt, err := s.db.Prepare("select id, idstr, name, mtime from list where idstr = ?", idstr)
	if !CheckSQLError(err) {
		return nil
	}

	err = stmt.Step()
	if !CheckSQLError(err) {
		return nil
	}

	values := stmt.Row()
	if values == nil || len(values) != 4 || values[0] == nil {
		// invalid/non-existant list
		return nil
	}

	id := values[0].(int64)
	l.Id = values[1].(string)
	l.Name = values[2].(string)
	l.ModificationTime = values[3].(int64)

	// select all the items, sorting by ID, since the parent items
       	// always have a smaller ID than their childrens
	stmt, err = s.db.Prepare("select parent_id, id, name, descr, status from item where list_id = ? order by id", id)
	if !CheckSQLError(err) {
		return nil
	}

	_, err = stmt.All(func(st *sqlite3.Statement, values ...interface{}) {
		parentId := values[0].(int64)
		item := &TodoListItem{
			values[1].(int64),
			values[2].(string),
			values[3].(string),
			values[4].(string),
			nil}
		if (parentId < 0) {
			l.AddItem(item)
		} else {
			parent := l.GetItem(parentId)
			if parent == nil {
				log.Println("Error when getting an item of list %s, ignoring item", l.Id);
				return
			}
			parent.AddItem(item)
		}
	})

	if !CheckSQLError(err) {
		return nil
	}

	return l
}

func (s *SQLiteStore) Exists(id string) bool {
	stmt, err := s.db.Prepare("select count(*) from list where idstr = ?", id)
	if !CheckSQLError(err) {
		return false
	}

	err = stmt.Step()
	if !CheckSQLError(err) {
		return false
	}

	n := stmt.Row()[0].(int64)
	stmt.Finalize()

	return n != 0
}

func (s *SQLiteStore) Delete(id string) bool {
	// this will also delete all the items thanks to the sql constraints
	stmt, err := s.db.Prepare("delete from list where idstr = ?", id)
	if !CheckSQLError(err) {
		return false
	}

	err = stmt.Step()
	if !CheckSQLError(err) {
		return false
	}
	return true
}

func (s *SQLiteStore) Set(list *TodoList) bool {
	_, err := s.db.Execute("begin transaction")
	if !CheckSQLError(err) {
		return false
	}

	if s.Exists(list.Id) {
		// To update, delete it first (maybe not the best solution)
		if !s.Delete(list.Id) {
			return false
		}
	}

	stmt, err := s.db.Prepare("insert into list(idstr, name, mtime) values (?, ?, ?)",
		list.Id, list.Name, list.ModificationTime)
	if !CheckSQLError(err) {
		s.db.Execute("end transaction")
		return false
	}

	err = stmt.Step()
	if !CheckSQLError(err) {
		s.db.Execute("end transaction")
		return false
	}

	id := s.db.LastInsertRowID()

	for _, item := range list.Items {
		stmt, err := s.db.Prepare("insert into item(parent_id, list_id, name, descr, status) values (-1, ?, ?, ?, ?)",
			id, item.Name, item.Description, item.Status)
		if !CheckSQLError(err) {
			s.db.Execute("end transaction")
			return false
		}

		err = stmt.Step()
		if !CheckSQLError(err) {
			s.db.Execute("end transaction")
			return false
		}

		item.Id = s.db.LastInsertRowID()

		if !s.AddItems(item, id) {
			s.db.Execute("end transaction")
			return false;
		}
	}

	_, err = s.db.Execute("end transaction")
	if !CheckSQLError(err) {
		return false
	}

	return true
}

func (s *SQLiteStore) AddItems(i *TodoListItem, id int64) bool {
	for _, item := range i.Items {
		stmt, err := s.db.Prepare("insert into item(parent_id, list_id, name, descr, status) values (?, ?, ?, ?, ?)",
			i.Id, id, item.Name, item.Description, item.Status)
		if !CheckSQLError(err) {
			return false
		}

		err = stmt.Step()
		if !CheckSQLError(err) {
			return false
		}

		item.Id = s.db.LastInsertRowID()

		if !s.AddItems(item, id) {
			return false
		}
	}
	return true
}