package main

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	db Database
}

func (h *Handler) insertNote(c *gin.Context) {
	var note Note

	c.BindJSON(&note)
	h.db.dataInsert(note)
	c.JSON(200, gin.H{
		"message": "Hello World",
		"name":    note.Name,
	})
}
