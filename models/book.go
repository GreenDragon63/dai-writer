package models

import (
	"dai-writer/auth"
	"strconv"
)

const prefixBook string = "Books/"

type Book struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Scenes      []int  `json:"scenes" binding:"dive"`
}

func (b *Book) setId(id int) {
	b.Id = id
}

func ListBook(u *auth.User) ([]*Book, bool) {
	path := prefixBook + strconv.Itoa(u.Id) + "/"
	return listJson[*Book](path, u.Id)
}

func LoadBook(u *auth.User, id int) (*Book, bool) {
	path := prefixBook + strconv.Itoa(u.Id) + "/"
	return loadJson[*Book](path, u.Id, id)
}

func SaveBook(u *auth.User, id int, postData Book) bool {
	path := prefixBook + strconv.Itoa(u.Id) + "/"
	return saveJson[*Book](path, u.Id, id, &postData)
}

func DeleteBook(u *auth.User, id int) bool {
	path := prefixBook + strconv.Itoa(u.Id) + "/"
	return deleteJson(path, u.Id, id)
}
