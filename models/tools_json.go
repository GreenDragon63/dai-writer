package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Identifiable interface {
	setId(id int)
}

func getId(prefix string) int {
	var id int = 0

	path := prefix
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

func listJson[T Identifiable](prefix string, uid int) ([]T, bool) {
	var data []T

	um.lock(uid)
	defer um.unlock(uid)

	path := prefix
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err.Error())
		return data, false
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err.Error())
		return data, false
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			content, err := ioutil.ReadFile(path + file.Name())
			if err != nil {
				log.Println(err.Error())
				return data, false
			}

			var item T
			err = json.Unmarshal(content, &item)
			if err != nil {
				log.Println(err.Error())
				return data, false
			}
			data = append(data, item)
		}
	}

	return data, true
}

func loadJson[T Identifiable](prefix string, uid int, id int) (T, bool) {
	var data T

	um.lock(uid)
	defer um.unlock(uid)

	path := prefix + strconv.Itoa(id) + ".json"

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		return data, false
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Println(err.Error())
		return data, false
	}

	return data, true
}

func saveJson[T Identifiable](prefix string, uid int, id int, data T) bool {
	var final_id int = 0
	var path string

	um.lock(uid)
	defer um.unlock(uid)

	if id == 0 {
		final_id = getId(prefix)
	} else {
		final_id = id
		data.setId(final_id)
	}

	if final_id == 0 {
		return false
	}

	path = prefix
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	path = prefix + "/" + strconv.Itoa(final_id) + ".json"

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func deleteJson(prefix string, uid int, id int) bool {
	um.lock(uid)
	defer um.unlock(uid)

	path := prefix + "/" + strconv.Itoa(id) + ".json"
	salvage := prefix + "/" + strconv.Itoa(id) + ".json.del"
	err := os.Rename(path, salvage)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
