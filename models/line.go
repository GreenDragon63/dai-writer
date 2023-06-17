package models

import (
	"dai-writer/auth"
	"strconv"
)

const prefixLine string = "Lines/"

type Line struct {
	Id          int      `json:"id"`
	BookId      int      `json:"book_id"  binding:"required"`
	SceneId     int      `json:"scene_id"  binding:"required"`
	Displayed   bool     `json:"displayed"`
	CharacterId int      `json:"character_id" binding:"required"`
	Content     []string `json:"content" binding:"dive"`
	Current     int      `json:"current"`
	Tokens      int      `json:"token"`
}

func (l *Line) setId(id int) {
	l.Id = id
}

func ListLine(u *auth.User, book int, scene int) ([]*Line, bool) {
	path := prefixLine + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return listJson[*Line](path, u.Id)
}

func LoadLine(u *auth.User, book int, scene int, id int) (*Line, bool) {
	path := prefixLine + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return loadJson[*Line](path, u.Id, id)
}

func SaveLine(u *auth.User, book int, scene int, id int, postData Line) bool {
	path := prefixLine + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return saveJson[*Line](path, u.Id, id, &postData)
}

func DeleteLine(u *auth.User, book int, scene int, id int) bool {
	path := prefixLine + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/" + strconv.Itoa(scene) + "/"
	return deleteJson(path, u.Id, id)
}
