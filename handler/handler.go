package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zashib/gin_server/entities"
)

func InsertNote(c *gin.Context) {
	var note entities.Note

	c.BindJSON(&note)
	h.Db.InsertNote(note)
	c.JSON(200, gin.H{
		"message": "Hello World",
		"title":   note.Title,
	})
}

func UpdateNote(c *gin.Context) {
	var note db.Note

	c.BindJSON(&note)
	h.Db.UpdateNote("Yes", "blabla")
	c.JSON(200, gin.H{
		"message": "Hello World",
		"title":   note.Title,
	})
}
