package models

import (
	"dai-writer/auth"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

const prefixCharacter string = "Characters/"

type Character struct {
	Id          int    `json:"id"`
	IsHuman     bool   `json:"is_human"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Personality string `json:"personality"`
	FirstMes    string `json:"first_mes"`
	MesExample  string `json:"mes_example"`
	Scenario    string `json:"scenario"`
}

type CharacterV2 struct {
	Data Character `json:"data"`
}

type CharacterInfos struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Scenario    string `json:"scenario"`
}

func (c *Character) setId(id int) {
	c.Id = id
}

func ListCharacterInfos(u *auth.User) ([]*CharacterInfos, bool) {
	var list []*Character
	var result []*CharacterInfos
	var ok bool
	path := prefixCharacter + strconv.Itoa(u.Id) + "/"
	list, ok = listJson[*Character](path, u.Id)
	if !ok {
		return result, false
	}
	for _, v := range list {
		result = append(result, &CharacterInfos{Id: v.Id, Name: v.Name, Description: v.Description, Scenario: v.Scenario})
	}
	return result, true
}

func ListCharacter(u *auth.User) ([]*Character, bool) {
	path := prefixCharacter + strconv.Itoa(u.Id) + "/"
	return listJson[*Character](path, u.Id)
}

func LoadCharacter(u *auth.User, id int) (*Character, bool) {
	path := prefixCharacter + strconv.Itoa(u.Id) + "/"
	return loadJson[*Character](path, u.Id, id)
}

func SaveCharacter(u *auth.User, id int, postData Character) bool {
	path := prefixCharacter + strconv.Itoa(u.Id) + "/"
	return saveJson[*Character](path, u.Id, id, &postData)
}

func DeleteCharacter(u *auth.User, id int) bool {
	path := prefixCharacter + strconv.Itoa(u.Id) + "/"
	return deleteJson(path, u.Id, id)
}

func UploadCharacterPath(u *auth.User) string {
	var id int = 0

	path := prefixCharacter + strconv.Itoa(u.Id) + "/"
	id = getId(path)
	if id == 0 {
		return ""
	}
	return prefixCharacter + strconv.Itoa(u.Id) + "/" + strconv.Itoa(id) + ".png"
}

func AvatarCharacterPath(u *auth.User, id int) string {
	return prefixCharacter + strconv.Itoa(u.Id) + "/" + strconv.Itoa(id) + ".png"
}

func DecodeCharacter(fileName string) bool {
	var chara Character
	var charav2 CharacterV2

	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer file.Close()

	header := make([]byte, 8)
	_, err = io.ReadFull(file, header)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	expectedHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytesEqual(header, expectedHeader) {
		log.Println("Not a PNG file : " + fileName)
		return false
	}

	for {
		sizeBytes := make([]byte, 4)
		_, err = io.ReadFull(file, sizeBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Cannot read chunk size: " + fileName)
			os.Remove(fileName)
			return false
		}
		chunkSize := binary.BigEndian.Uint32(sizeBytes)

		nameBytes := make([]byte, 4)
		_, err = io.ReadFull(file, nameBytes)
		if err != nil {
			log.Println("Cannot read chunk name: " + fileName)
			os.Remove(fileName)
			return false
		}
		chunkName := string(nameBytes)

		chunkData := make([]byte, chunkSize)
		_, err = io.ReadFull(file, chunkData)
		if err != nil {
			log.Println("Cannot read chunk data: " + fileName)
			os.Remove(fileName)
			return false
		}
		// Do not check the CRC
		_, err = file.Seek(4, io.SeekCurrent)
		if err != nil {
			log.Println("Cannot read CRC: " + fileName)
			os.Remove(fileName)
			return false
		}

		if chunkName == "tEXt" && string(chunkData[:5]) == "chara" {
			charaBase64 := string(chunkData[6:])
			decodedChara, err := base64.StdEncoding.DecodeString(charaBase64)
			if err != nil {
				log.Println("chara chunk founded, but it is not encoded in base64: " + fileName)
				os.Remove(fileName)
				return false
			}
			err = json.Unmarshal(decodedChara, &chara)
			if err != nil {
				log.Println("chara chunk founded, but it is not with the good format: " + fileName)
				os.Remove(fileName)
				return false
			}
			if chara.Name == "" {
				err = json.Unmarshal(decodedChara, &charav2)
				if err != nil {
					log.Println("chara chunk founded, but it is not with the good format: " + fileName)
					os.Remove(fileName)
					return false
				}
				chara = charav2.Data
			}
			id, err := strconv.Atoi(path.Base(fileName[:len(fileName)-len(".png")]))
			if err != nil {
				log.Println("Weird error while decoding the id: " + fileName)
				os.Remove(fileName)
				return false
			}
			chara.Id = id
			jsonData, err := json.Marshal(chara)
			if err != nil {
				log.Println("Weird error while recoding the character: " + fileName)
				return false
			}
			jsonFile := fileName[:len(fileName)-len(".png")] + ".json"
			err = os.WriteFile(jsonFile, jsonData, 0644)
			if err != nil {
				log.Println(err.Error())
				os.Remove(fileName)
				return false
			}
			return true
		}
	}
	os.Remove(fileName)
	return false
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
