package engram

import (
	"dai-writer/auth"
	"dai-writer/models"

	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type exportLineRequest struct {
	Content string `json:"content"`
	Present []int  `json:"present"`
}

func ExportLine(u *auth.User, line *models.Line) bool {
	var request *http.Request
	var response *http.Response
	var err error
	var requestData *exportLineRequest
	var api, key, url string
	var scene *models.Scene
	var ok bool

	scene, ok = models.LoadScene(u, line.BookId, line.SceneId)
	if ok != true {
		return false
	}
	requestData = &exportLineRequest{
		Content: line.Content[line.Current],
		Present: scene.Characters,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}

	api = os.Getenv("ENGRAM_API")
	url = api + "line/" + strconv.Itoa(u.Id) + "/" + strconv.Itoa(line.BookId) + "/" + strconv.Itoa(line.SceneId) + "/" + strconv.Itoa(line.Id)
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	key = os.Getenv("ENGRAM_API_KEY")
	request.Header.Add("API-Key", key)
	request.Header.Add("Content-Type", "application/json")

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return false
	}
	if response.StatusCode != 200 {
		return false
	}
	return true
}

func RemoveLine(u *auth.User, line *models.Line) bool {
	var request *http.Request
	var response *http.Response
	var err error
	var api, key, url string

	api = os.Getenv("ENGRAM_API")
	url = api + "line/" + strconv.Itoa(u.Id) + "/" + strconv.Itoa(line.BookId) + "/" + strconv.Itoa(line.SceneId) + "/" + strconv.Itoa(line.Id)
	request, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false
	}
	key = os.Getenv("ENGRAM_API_KEY")
	request.Header.Add("API-Key", key)

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return false
	}
	if response.StatusCode != 200 {
		return false
	}
	return true
}

func Search(u *auth.User, bookId, sceneId, lineId, charId int) ([]string, bool) {
	var request *http.Request
	var response *http.Response
	var err error
	var api, key, url string
	var result []string

	api = os.Getenv("ENGRAM_API")
	url = api + "search/" + strconv.Itoa(u.Id) + "/" + strconv.Itoa(bookId) + "/" + strconv.Itoa(sceneId) + "/" + strconv.Itoa(lineId) + "/" + strconv.Itoa(charId)
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return result, false
	}
	key = os.Getenv("ENGRAM_API_KEY")
	request.Header.Add("API-Key", key)

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return result, false
	}
	if response.StatusCode != 200 {
		return result, false
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, false
	}
	return result, true
}
