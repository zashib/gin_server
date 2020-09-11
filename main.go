package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/note", func(c *gin.Context) {
		body := struct {
			Name string
		}{}
		c.BindJSON(&body)
		c.JSON(200, gin.H{
			"message": "Hello World",
			"name":    body.Name,
		})
	})
	router.GET("/note/:name")
	router.PUT("/note/:name")
	router.DELETE("/note/:name")

	router.Run()
}
