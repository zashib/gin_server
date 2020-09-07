package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"xorm.io/xorm"
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5434", "user", "example", "pwd")
}

type Note struct {
	Id      int64
	Name    string
	Content string
	Status  time.Time
}

func dataInsert(note Note, engine *xorm.Engine) {
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
	_, err := engine.In("name", name).Update(&Note{Content: newContent})
	if err != nil {
		fmt.Println(err)
	}
}

func dataDeleteByName(name string, engine *xorm.Engine){
	_, err := engine.Delete(Note{Name: name})
	if err != nil {
		fmt.Println(err)
	}
}

// TODO autoupdate time func when data changed
// TODO set unique not null(default) values for cols
func main() {
	engine, err := xorm.NewEngine("postgres", connectionString())
	if err != nil {
		fmt.Println(err)
	}
	err = engine.Sync2(new(Note))
	if err != nil {
		fmt.Println(err)
	}

	note := Note{
		Name:    "test",
		Content: "some contetnt",
	}

	dataInsert(note, engine)

	name := "test1"
	fmt.Println(dataGetByName(name, engine))

	dataUpdateContentByName(name, "else", engine)

	dataDeleteByName(name, engine)

}
