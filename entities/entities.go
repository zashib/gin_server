package entities

import "time"

type Note struct {
	Id      int64
	Title   string `xorm:"notnull unique"`
	Content string
	Status  time.Time
}
