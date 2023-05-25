package models

import (
	"dai-writer/auth"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Character struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Personality string `json:"personality"`
	First_mes   string `json:"first_mes"`
	Mes_example string `json:"mes_example"`
	Scenario    string `json:"scenario"`
}

func LoadCharacter(u *auth.User, id int) (Character, error) {
	var chara Character

	path := "characters/" + strconv.Itoa(u.Id) + "/" + strconv.Itoa(id) + ".json"
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return chara, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&chara)
	if err != nil {
		log.Println(err.Error())
		return chara, err
	}
	return chara, nil
}

func UploadCharacterPath(u *auth.User) string {
	var id int = 0

	id = GetId("characters/", u.Id)
	if id == 0 {
		return ""
	}
	return "characters/" + strconv.Itoa(u.Id) + "/" + strconv.Itoa(id) + ".png"
}

func DecodeCharacter(u *auth.User, id int) {
	// TODO
}
