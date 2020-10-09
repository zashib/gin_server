package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"strconv"

	"github.com/gin-gonic/gin"
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

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5434", "user", "example", "pwd")
}

// TODO:

// Initialize xorm db and put in context
// Implement all methods for notes/tasks and categories
// Write api tests

// https://github.com/gavv/httpexpect
// Test: get tasks -> toggle_task -> get tasks -> поверяется что bool у объекта поменялся
// Test: get tasks -> change name -> get tasks -> поверяется что у объекта поменялся

func main() {
	engine, engineErr := xorm.NewEngine("postgres", connectionString())
	if engineErr != nil {
		fmt.Println(engineErr)
	}

	syncErr := engine.Sync2(new(Tasks), new(Category))
	if syncErr != nil {
		fmt.Println(syncErr)
	}

	router := gin.Default()

	// Создаёт задачу
	router.POST("/task", func(c *gin.Context) {
		var task Tasks

		bindErr := c.BindJSON(&task)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, insertErr := engine.Insert(&task)
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
	// Возвращает все задачи
	router.GET("/task", func(c *gin.Context) {
		type TasksCategory struct {
			Title   string
			Content string
			Status  bool
			Name    string
		}
		var tasks []TasksCategory
		findErr := engine.Table("tasks").Join("LEFT", "category", "category.id = tasks.category_id").
			Find(&tasks)
		if findErr != nil {
			fmt.Println(findErr)
		}
		c.JSON(200, gin.H{
			"tasks": tasks,
		})
	})
	// Переключает задачу невыполнено/выполнено
	router.PUT("/toggle_task/:id", func(c *gin.Context) {
		id := c.Param("id")

		var valuesMap = make(map[string]string)
		has, getErr := engine.Table(&Tasks{}).Where("id = ?", id).Get(&valuesMap)
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

		_, updateErr := engine.UseBool("status").ID(id).Update(&Tasks{
			Status: !parseBool,
		})

		if updateErr != nil {
			fmt.Println(updateErr)
		}
	})

	// Меняет контент и название задачи
	router.PUT("/task/:id", func(c *gin.Context) {
		id := c.Param("id")

		var newTask Tasks

		bindErr := c.BindJSON(&newTask)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, updateErr := engine.ID(id).Update(&Tasks{
			Title:   newTask.Title,
			Content: newTask.Content,
		})

		if updateErr != nil {
			fmt.Println(updateErr)
		}
	})
	// Удаляет
	router.DELETE("/task/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, deleteErr := engine.ID(id).Delete(&Tasks{})

		if deleteErr != nil {
			fmt.Println(deleteErr)
		}
	})

	// Создаёт категорию(странное поведение при добавлении нового значения)
	router.POST("/category", func(c *gin.Context) {
		var category Category

		bindErr := c.BindJSON(&category)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, insertErr := engine.Insert(&category)
		if insertErr != nil {
			fmt.Println(insertErr)
		}

		c.JSON(200, gin.H{
			"id":   category.Id,
			"name": category.Name,
		})
	})
	// Возвращает все категории
	router.GET("/category", func(c *gin.Context) {
		var categories []Category
		findErr := engine.Find(&categories)
		fmt.Println(categories)
		if findErr != nil {
			fmt.Println(findErr)
		}

		c.JSON(200, gin.H{
			"categories": categories,
		})
	})

	// Меняет имя категории
	router.PUT("/category/:id", func(c *gin.Context) {
		id := c.Param("id")

		var category Category

		bindErr := c.BindJSON(&category)
		if bindErr != nil {
			fmt.Println(bindErr)
		}

		_, updateErr := engine.ID(id).Update(&Category{
			Name: category.Name,
		})

		if updateErr != nil {
			fmt.Println(updateErr)
		}
	})

	// Удаляет категорию
	router.DELETE("/category/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, deleteErr := engine.ID(id).Delete(&Category{})

		if deleteErr != nil {
			fmt.Println(deleteErr)
		}
	})

	runErr := router.Run()
	if runErr != nil {
		fmt.Println(runErr)
	}
}
