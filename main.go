package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zashib/gin_server/handler"
)

// TODO:

// Initialize xorm db and put in context
// Implement all methods for notes/tasks and categories
// Write api tests

// https://github.com/gavv/httpexpect
// Test: get tasks -> toggle_task -> get tasks -> поверяется что bool у объекта поменялся
// Test: get tasks -> change name -> get tasks -> поверяется что у объекта поменялся

func main() {
	router := gin.Default()

	// Создаёт задачу
	router.POST("/task", handler.InsertNote)
	// Возвращает все задачи
	router.GET("/task")
	// Переключает задачу невыполнено/выполнено
	router.PUT("/toggle_task/:id")
	// Меняет контент и название задачи
	router.PUT("/task/:id")
	// Удаляет
	router.DELETE("/task/:id")

	// Создаёт капегорию
	router.POST("/category")
	// Возвращает все категории
	router.GET("/category/:id")
	// Меняет имя категории body {"name":"bla"}
	router.PUT("/category/:id")
	// Удаляет
	router.DELETE("/category/:id")

	router.Run()
}
