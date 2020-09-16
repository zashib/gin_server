package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zashib/gin_server/db"
)

type Handler struct {
	Db db.Database
}

func (h *Handler) InsertNote(c *gin.Context) {
	var note db.Note

	c.BindJSON(&note)
	h.Db.DataInsert(note)
	c.JSON(200, gin.H{
		"message": "Hello World",
		"name":    note.Name,
	})
}
