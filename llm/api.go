package llm

import (
	"bytes"
	"dai-writer/aes"
	"encoding/json"
	"io/ioutil"
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

type CompletionRequest struct {
	Prompt           string  `json:"prompt"`
	UseStory         bool    `json:"use_story"`
	UseMemory        bool    `json:"use_memory"`
	UseAuthorsNote   bool    `json:"use_authors_note"`
	UseWorldInfo     bool    `json:"use_world_info"`
	MaxContextLength int     `json:"max_context_length"`
	MaxLength        int     `json:"max_length"`
	RepPen           float64 `json:"rep_pen"`
	RepPenRange      int     `json:"rep_pen_range"`
	RepPenSlope      float64 `json:"rep_pen_slope"`
	Temperature      float64 `json:"temperature"`
	Tfs              float64 `json:"tfs"`
	TopA             int     `json:"top_a"`
	TopK             int     `json:"top_k"`
	TopP             float64 `json:"top_p"`
	Typical          int     `json:"typical"`
	SamplerOrder     []int   `json:"sampler_order"`
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

	body, err := ioutil.ReadAll(response.Body)
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

func GetCompletion(text string) (string, bool) {
	var url, message, key, result string
	var tokens int
	var requestData CompletionRequest

	key = os.Getenv("PSK")
	if key != "" {
		message = aes.StrEncrypt(text, key)
	} else {
		message = text
	}

	requestData = CompletionRequest{
		Prompt:           message,
		UseStory:         false,
		UseMemory:        false,
		UseAuthorsNote:   false,
		UseWorldInfo:     false,
		MaxContextLength: 2048,
		MaxLength:        50,
		RepPen:           1.1,
		RepPenRange:      1024,
		RepPenSlope:      0.9,
		Temperature:      0.65,
		Tfs:              0.9,
		TopA:             0,
		TopK:             0,
		TopP:             0.9,
		Typical:          1,
		SamplerOrder:     []int{6, 0, 1, 2, 3, 4, 5},
	}

	if key != "" {
		url = ApiUrl() + "api/v1.1/generate"
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

	body, err := ioutil.ReadAll(response.Body)
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
