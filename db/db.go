package db

import (
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"xorm.io/xorm"
)

type Note struct {
	Id      int64
	Title   string `xorm:"notnull unique"`
	Content string
	Status  time.Time
}

type Database struct {
	engine xorm.Engine
}

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5434", "user", "example", "pwd")
}

func (db Database) Init() Database {
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

func (db Database) InsertNote(note Note) {
	note.Status = time.Now()
	_, err := db.engine.InsertOne(&note)
	if err != nil {
		fmt.Println(err)
	}
}

func (db Database) GetNoteById(title string) interface{} {
	var note = Note{
		Title: title,
	}
	_, err := db.engine.Get(&note)
	if err != nil {
		fmt.Println(err)
	}
	return note
}

func (db Database) UpdateNote(title string, newContent string) {
	_, err := db.engine.In("title", title).Update(&Note{Content: newContent, Status: time.Now()})
	if err != nil {
		fmt.Println(err)
	}
}

func (db Database) DeleteNote(title string) {
	_, err := db.engine.Delete(Note{Title: title})
	if err != nil {
		fmt.Println(err)
	}
}

// TODO change CRUD methods logic from by 'name' to by 'id'
// TODO set unique not null(default) values for cols
//func main() {
//	db := Database{}.init()
//
//	note := Note{
//		Title:    "test",
//		Content: "some contetnt",
//	}
//
//	db.dataInsert(note)
//
//	name := "test1"
//	fmt.Println(db.GetNoteById(name))
//
//	db.UpdateNote(name, "else")
//
//	name = "test2"
//	db.DeleteNote(name)
//
//}
