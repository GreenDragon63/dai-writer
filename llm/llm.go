package llm

import (
	"bytes"
	"dai-writer/auth"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TokenRequest struct {
	Prompt string `json:"prompt"`
}

type TokenResponse struct {
	Results []struct {
		Tokens int `json:"tokens"`
	} `json:"results"`
}

func Generate(u *auth.User, book int, scene int, id int) string {
	return "Generated text"
}

func ApiUrl() string {
	return os.Getenv("BACKEND_API")
}

func GetTokens(text string) int {
	var url string
	var requestData TokenRequest

	url = ApiUrl() + "api/v1/token-count"
	requestData.Prompt = text

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Failed to convert JSON data : %s\n", err.Error())
		return 0
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Can't connect to backend : %s\n", err.Error())
		return 0
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Can't read the response : %s\n", err.Error())
		return 0
	}

	var data TokenResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("can't decode JSON : %s\n", err.Error())
		return 0
	}

	return data.Results[0].Tokens
}
