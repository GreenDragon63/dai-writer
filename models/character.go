package models

import (
	"dai-writer/auth"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
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
	path := "characters/" + strconv.Itoa(u.Id) + "/"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	for _, file := range files {
		f := strings.Split(file.Name(), ".")
		i, err := strconv.Atoi(f[0])
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		if i > id {
			id = i
		}
	}
	id++
	return "characters/" + strconv.Itoa(u.Id) + "/" + strconv.Itoa(id) + ".png"
}

func DecodeCharacter(u *auth.User, id int) {
	// TODO
}
