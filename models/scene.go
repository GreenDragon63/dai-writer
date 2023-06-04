package models

import (
	"dai-writer/auth"
	"strconv"
)

const prefixScene string = "Scenes/"

type Scene struct {
	Id             int    `json:"id"`
	Description    string `json:"description" binding:"required"`
	Characters     []int  `json:"characters" binding:"required,dive"`
	First_line     int    `json:"first_line"`
	Previous_scene int    `json:"previous_scene" binding:"required"`
	Next_scene     int    `json:"next_scene" binding:"required"`
}

func (s *Scene) setId(id int) {
	s.Id = id
}

func ListScene(u *auth.User, book int) ([]*Scene, bool) {
	path := prefixScene + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/"
	return listJson[*Scene](path, u.Id)
}

func LoadScene(u *auth.User, book int, id int) (*Scene, bool) {
	path := prefixScene + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/"
	return loadJson[*Scene](path, u.Id, id)
}

func SaveScene(u *auth.User, book int, id int, postData Scene) bool {
	path := prefixScene + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/"
	return saveJson[*Scene](path, u.Id, id, &postData)
}

func DeleteScene(u *auth.User, book int, id int) bool {
	path := prefixScene + strconv.Itoa(u.Id) + "/" + strconv.Itoa(book) + "/"
	return deleteJson(path, u.Id, id)
}
