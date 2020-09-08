package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/note")
	router.GET("/note/:name")
	router.PUT("/note/:name")
	router.DELETE("/note/:name")

	router.Run()
}
