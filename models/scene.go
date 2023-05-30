package models

import (
	"dai-writer/auth"
)

const prefixScene string = "Scenes/"

type Scene struct {
	Description    string `json:"description" binding:"required"`
	Characters     []int  `json:"characters" binding:"required,dive"`
	First_line     int    `json:"first_line"`
	Previous_scene int    `json:"previous_scene" binding:"required"`
	Next_scene     int    `json:"next_scene" binding:"required"`
}

func ListScene(u *auth.User) ([]Scene, bool) {
	return listJson[Scene](prefixScene, u.Id)
}

func LoadScene(u *auth.User, id int) (Scene, bool) {
	return loadJson[Scene](prefixScene, u.Id, id)
}

func SaveScene(u *auth.User, id int, postData Scene) bool {
	return saveJson[Scene](prefixScene, u.Id, id, postData)
}

func DeleteScene(u *auth.User, id int) bool {
	return deleteJson(prefixScene, u.Id, id)
}
