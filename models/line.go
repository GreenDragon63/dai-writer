package models

import (
	"dai-writer/auth"
)

const prefixLine string = "Lines/"

type Line struct {
	Character     int    `json:"character" binding:"required"`
	Content       string `json:"content"`
	Previous_Line int    `json:"previous_Line" binding:"required"`
	Next_Line     int    `json:"next_Line" binding:"required"`
}

func ListLine(u *auth.User) ([]Line, bool) {
	return listJson[Line](prefixLine, u.Id)
}

func LoadLine(u *auth.User, id int) (Line, bool) {
	return loadJson[Line](prefixLine, u.Id, id)
}

func SaveLine(u *auth.User, id int, data []byte) bool {
	return saveJson(prefixLine, u.Id, id, data)
}

func DeleteLine(u *auth.User, id int) bool {
	return deleteJson(prefixLine, u.Id, id)
}
