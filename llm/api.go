package llm

import (
	"bytes"
	"dai-writer/aes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type ModelResponse struct {
	Result string `json:"result"`
}

type TokenRequest struct {
	Prompt string `json:"prompt"`
}

type TokenResponse struct {
	Results []struct {
		Tokens int `json:"tokens"`
	} `json:"results"`
}

type CompletionResponse struct {
	Results []struct {
		Text string `json:"text"`
	} `json:"results"`
}

func ApiUrl() string {
	return os.Getenv("BACKEND_API")
}

func GetModel() string {
	var url string

	url = ApiUrl() + "api/v1/model"

	response, err := http.Get(url)
	if err != nil {
		log.Printf("Can't connect to backend : %s\n", err.Error())
		return ""
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Can't read the response : %s\n", err.Error())
		return ""
	}

	var data ModelResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("can't decode JSON : %s\n", err.Error())
		return ""
	}

	return data.Result
}

func GetTokens(text string) int {
	var url, key string
	var requestData TokenRequest

	key = os.Getenv("PSK")
	if key != "" {
		requestData.Prompt = aes.StrEncrypt(text, key)
	} else {
		requestData.Prompt = text
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Failed to convert JSON data : %s\n", err.Error())
		return 0
	}

	if key != "" {
		url = ApiUrl() + "api/v2/token-count"
	} else {
		url = ApiUrl() + "api/v1/token-count"
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Can't connect to backend : %s\n", err.Error())
		return 0
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
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

// returned bool is true if the completion is finished
func GetCompletion(text string) (string, bool) {
	var model, url, message, key, result string
	var tokens int
	var requestData *Model
	var err error

	key = os.Getenv("PSK")
	if key != "" {
		message = aes.StrEncrypt(text, key)
	} else {
		message = text
	}

	model = GetModel()
	requestData, err = loadModelConfig(model)
	if err != nil {
		log.Printf("Can't load model config : %s\n", err.Error())
		return "", true
	}
	requestData.Prompt = message
	requestData.MaxNewTokens = 50

	if key != "" {
		url = ApiUrl() + "api/v2/generate"
	} else {
		url = ApiUrl() + "api/v1/generate"
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Failed to convert JSON data : %s\n", err.Error())
		return "", true
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Can't connect to backend : %s\n", err.Error())
		return "", true
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Can't read the response : %s\n", err.Error())
		return "", true
	}

	var data CompletionResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("can't decode JSON : %s\n", err.Error())
		return "", true
	}

	if key != "" {
		result = aes.StrDecrypt(data.Results[0].Text, key)
	} else {
		result = data.Results[0].Text
	}

	tokens = GetTokens(result)
	log.Printf("Tokens : %d\n", tokens)
	if tokens < 50 {
		return result, true
	} else {
		return result, false
	}
}
