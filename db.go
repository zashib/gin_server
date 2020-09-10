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

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5434", "user", "example", "pwd")
}

func dataInsert(note Note, engine *xorm.Engine) {
	note.Status = time.Now()
	_, err := engine.InsertOne(&note)
	if err != nil {
		fmt.Println(err)
	}
}

func dataGetByName(name string, engine *xorm.Engine) interface{} {
	var note = Note{
		Name: name,
	}
	_, err := engine.Get(&note)
	if err != nil {
		fmt.Println(err)
	}
	return note
}

func dataUpdateContentByName(name string, newContent string, engine *xorm.Engine) {
	_, err := engine.In("name", name).Update(&Note{Content: newContent, Status: time.Now()})
	if err != nil {
		fmt.Println(err)
	}
}

func dataDeleteByName(name string, engine *xorm.Engine) {
	_, err := engine.Delete(Note{Name: name})
	if err != nil {
		fmt.Println(err)
	}
}

func createDB() *xorm.Engine {
	engine, err := xorm.NewEngine("postgres", connectionString())
	if err != nil {
		fmt.Println(err)
	}
	err = engine.Sync2(new(Note))
	if err != nil {
		fmt.Println(err)
	}
	return engine
}

// TODO autoupdate time func when data changed
// TODO set unique not null(default) values for cols
func main() {
	engine := createDB()

	note := Note{
		Name:    "test",
		Content: "some contetnt",
	}

	dataInsert(note, engine)

	name := "test1"
	fmt.Println(dataGetByName(name, engine))

	dataUpdateContentByName(name, "else", engine)

	name = "test2"
	dataDeleteByName(name, engine)

}
