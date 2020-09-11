package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"xorm.io/xorm"
)

type Note struct {
	Id      int64
	Name    string `xorm:"notnull unique"`
	Content string
	Status  time.Time
}

type Database struct {
	engine xorm.Engine
}

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5434", "user", "example", "pwd")
}

func (db Database) init() Database {
	engine, err := xorm.NewEngine("postgres", connectionString())
	if err != nil {
		fmt.Println(err)
	}
	err = engine.Sync2(new(Note))
	if err != nil {
		fmt.Println(err)
	}
	db.engine = *engine
	return db
}

func (db Database) dataInsert(note Note) {
	note.Status = time.Now()
	_, err := db.engine.InsertOne(&note)
	if err != nil {
		fmt.Println(err)
	}
}

func (db Database) dataGetByName(name string) interface{} {
	var note = Note{
		Name: name,
	}
	_, err := db.engine.Get(&note)
	if err != nil {
		fmt.Println(err)
	}
	return note
}

func (db Database) dataUpdateContentByName(name string, newContent string) {
	_, err := db.engine.In("name", name).Update(&Note{Content: newContent, Status: time.Now()})
	if err != nil {
		fmt.Println(err)
	}
}

func (db Database) dataDeleteByName(name string) {
	_, err := db.engine.Delete(Note{Name: name})
	if err != nil {
		fmt.Println(err)
	}
}

// TODO set unique not null(default) values for cols
func main() {
	db := Database{}.init()

	note := Note{
		Name:    "test",
		Content: "some contetnt",
	}

	db.dataInsert(note)

	name := "test1"
	fmt.Println(db.dataGetByName(name))

	db.dataUpdateContentByName(name, "else")

	name = "test2"
	db.dataDeleteByName(name)

}
