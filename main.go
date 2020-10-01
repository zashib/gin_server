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
	router.GET("/note/:id")
	router.PUT("/note/:id")
	router.DELETE("/note/:id")

	router.Run()
}
