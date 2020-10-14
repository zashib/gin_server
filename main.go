package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"strconv"
	"xorm.io/xorm"
)

type Tasks struct {
	Id         int64
	Title      string `xorm:"notnull unique"`
	Content    string
	Status     bool  `xorm:"notnull"`
	CategoryId int64 `xorm:"'category_id' index"`
}

type Category struct {
	Id   int64
	Name string `xorm:"unique"`
}

type TasksCategory struct {
	Id      int64
	Title   string
	Content string
	Status  bool
	Name    string
}

var (
	r = gin.New()
)

// TODO:

// Initialize xorm initDb and put in context
// Implement all methods for notes/tasks and categories
// Write api tests

// https://github.com/gavv/httpexpect
// Test: get tasks -> toggle_task -> get tasks -> поверяется что bool у объекта поменялся
// Test: get tasks -> change name -> get tasks -> поверяется что у объекта поменялся

func main() {
	db := createDB(connectionString())
	// Создаёт задачу
	postTask(db)
	// Возвращает все задачи
	getTasks(db)
	// Переключает задачу невыполнено/выполнено
	toggleTask(db)
	// Меняет контент и название задачи
	updateTask(db)
	// Удаляет задачу
	r.DELETE("/task/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, deleteErr := db.ID(id).Delete(&Tasks{})

		if deleteErr != nil {
			fmt.Println(deleteErr)
		}
	})

	// Создаёт категорию(странное поведение при добавлении нового значения)
	r.POST("/category", func(c *gin.Context) {
		var category Category

		bindErr := c.BindJSON(&category)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, insertErr := db.Insert(&category)
		if insertErr != nil {
			fmt.Println(insertErr)
		}

		c.JSON(200, gin.H{
			"id":   category.Id,
			"name": category.Name,
		})
	})
	// Возвращает все категории
	r.GET("/category", func(c *gin.Context) {
		var categories []Category
		findErr := db.Find(&categories)
		fmt.Println(categories)
		if findErr != nil {
			fmt.Println(findErr)
		}

		c.JSON(200, gin.H{
			"categories": categories,
		})
	})

	// Меняет имя категории
	r.PUT("/category/:id", func(c *gin.Context) {
		id := c.Param("id")

		var category Category

		bindErr := c.BindJSON(&category)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, updateErr := db.ID(id).Update(&Category{
			Name: category.Name,
		})

		if updateErr != nil {
			fmt.Println(updateErr)
		}
	})

	// Удаляет категорию
	r.DELETE("/category/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, deleteErr := db.ID(id).Delete(&Category{})

		if deleteErr != nil {
			fmt.Println(deleteErr)
		}
	})

	runErr := r.Run()
	if runErr != nil {
		fmt.Println(runErr)
	}
}

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5434", "user", "example", "pwd")
}

func createDB(connection string) *xorm.Engine {
	db, engineErr := xorm.NewEngine("postgres", connection)
	if engineErr != nil {
		fmt.Println(engineErr)
	}

	syncErr := db.Sync2(new(Tasks), new(Category))
	if syncErr != nil {
		fmt.Println(syncErr)
	}
	return db
}

func postTask(db *xorm.Engine) *gin.Engine {
	r.POST("/task", func(c *gin.Context) {
		var task Tasks

		bindErr := c.BindJSON(&task)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, insertErr := db.Insert(&task)
		if insertErr != nil {
			fmt.Println(insertErr)
		}

		c.JSON(200, gin.H{
			"id":      task.Id,
			"title":   task.Title,
			"content": task.Content,
			"status":  task.Status,
		})
	})
	return r
}

func getTasks(db *xorm.Engine) *gin.Engine {
	r.GET("/task", func(c *gin.Context) {
		tasks := make(map[int64]TasksCategory)
		findErr := db.Table("tasks").Join("LEFT", "category", "category.id = tasks.category_id").
			Find(&tasks)
		if findErr != nil {
			fmt.Println(findErr)
		}
		c.JSON(200, gin.H{
			"tasks": tasks,
		})
	})
	return r
}

func toggleTask(db *xorm.Engine) *gin.Engine {
	r.PUT("/toggle_task/:id", func(c *gin.Context) {
		id := c.Param("id")

		var valuesMap = make(map[string]string)
		has, getErr := db.Table(&Tasks{}).Where("id = ?", id).Get(&valuesMap)
		if !has {
			fmt.Printf("task with id %s doesn`t exist\n", id)
		}
		if getErr != nil {
			fmt.Println(getErr)
		}
		parseBool, parseErr := strconv.ParseBool(valuesMap["status"])
		if parseErr != nil {
			fmt.Println(parseErr)
		}

		_, updateErr := db.UseBool("status").ID(id).Update(&Tasks{
			Status: !parseBool,
		})

		if updateErr != nil {
			fmt.Println(updateErr)
		}
	})
	return r
}

func updateTask(db *xorm.Engine) *gin.Engine {
	r.PUT("/task/:id", func(c *gin.Context) {
		id := c.Param("id")

		var newTask Tasks

		bindErr := c.BindJSON(&newTask)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, updateErr := db.ID(id).Update(&Tasks{
			Title:   newTask.Title,
			Content: newTask.Content,
		})

		if updateErr != nil {
			fmt.Println(updateErr)
		}
	})
	return r
}
