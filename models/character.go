package models

import (
	"dai-writer/auth"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
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

func DecodeCharacter(fileName string) bool {
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
			jsonFile := fileName[:len(fileName)-len(".png")] + ".json"
			err = os.WriteFile(jsonFile, decodedChara, 0644)
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
