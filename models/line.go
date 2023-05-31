package models

import (
	"dai-writer/auth"
	"strconv"
)

const prefixLine string = "Lines/"

type Line struct {
	Character     int    `json:"character" binding:"required"`
	Content       string `json:"content"`
	Previous_Line int    `json:"previous_Line" binding:"required"`
	Next_Line     int    `json:"next_Line" binding:"required"`
}

func ListLine(u *auth.User, book int, scene int) ([]Line, bool) {
	path := prefixLine + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return listJson[Line](path, u.Id)
}

func LoadLine(u *auth.User, book int, scene int, id int) (Line, bool) {
	path := prefixLine + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return loadJson[Line](path, u.Id, id)
}

func SaveLine(u *auth.User, book int, scene int, id int, postData Line) bool {
	path := prefixLine + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return saveJson[Line](path, u.Id, id, postData)
}

func DeleteLine(u *auth.User, book int, scene int, id int) bool {
	path := prefixLine + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return deleteJson(path, u.Id, id)
}
