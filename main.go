package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db := Database{}.init()
	handler := Handler{db: db}

	router.POST("/note", handler.insertNote)
	router.GET("/note/:name")
	router.PUT("/note/:name")
	router.DELETE("/note/:name")

	router.Run()
}
