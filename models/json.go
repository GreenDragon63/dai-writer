package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func getId(prefix string, uid int) int {
	var id int = 0

	path := prefix + strconv.Itoa(uid) + "/"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	for _, file := range files {
		f := strings.Split(file.Name(), ".")
		i, err := strconv.Atoi(f[0])
		if err != nil {
			log.Println(err.Error())
			return 0
		}
		if i > id {
			id = i
		}
	}
	id++

	return id
}

func loadJson[T any](prefix string, uid int, id int) (T, bool) {
	var data T

	path := prefix + strconv.Itoa(uid) + "/" + strconv.Itoa(id) + ".json"
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return data, false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		return data, false
	}
	return data, true
}

func saveJson(prefix string, uid int, id int, jsonData []byte) bool {
	var final_id int = 0
	var path string

	if id == 0 {
		final_id = getId(prefix, uid)
	} else {
		final_id = id
	}

	if final_id == 0 {
		return false
	}

	path = prefix + strconv.Itoa(uid) + "/"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	path = prefix + strconv.Itoa(uid) + "/" + strconv.Itoa(final_id) + ".json"

	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func deleteJson(prefix string, uid int, id int) bool {
	path := "characters/" + strconv.Itoa(uid) + "/" + strconv.Itoa(id) + ".json"
	salvage := "characters/" + strconv.Itoa(uid) + "/" + strconv.Itoa(id) + ".json.del"
	err := os.Rename(path, salvage)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
