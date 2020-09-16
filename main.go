package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zashib/gin_server/db"
	"github.com/zashib/gin_server/handler"
)

func main() {
	router := gin.Default()
	db := db.Database{}.Init()
	handler := handler.Handler{Db: db}

	router.POST("/note", handler.InsertNote)
	router.GET("/note/:name")
	router.PUT("/note/:name")
	router.DELETE("/note/:name")

	router.Run()
}
