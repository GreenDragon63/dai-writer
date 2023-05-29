package models

import (
	"dai-writer/auth"
)

const prefixBook string = "Books/"

type Book struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	First_scene int    `json:"first_scene"`
}

func LoadBook(u *auth.User, id int) (Book, bool) {
	return loadJson[Book](prefixBook, u.Id, id)
}

func SaveBook(u *auth.User, id int, data []byte) bool {
	return saveJson(prefixBook, u.Id, id, data)
}

func DeleteBook(u *auth.User, id int) bool {
	return deleteJson(prefixBook, u.Id, id)
}
