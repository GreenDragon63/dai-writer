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

func LoadScene(u *auth.User, id int) (Scene, bool) {
	return loadJson[Scene](prefixScene, u.Id, id)
}

func SaveScene(u *auth.User, id int, data []byte) bool {
	return saveJson(prefixScene, u.Id, id, data)
}

func DeleteScene(u *auth.User, id int) bool {
	return deleteJson(prefixScene, u.Id, id)
}
