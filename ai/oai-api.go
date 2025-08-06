package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	CHAT_ENDPOINT = "https://untamed.ai/v1/chat/completions"
	COMP_ENDPOINT = "https://untamed.ai/v1/completions"
	MODEL         = "llama"
	MODEL_CTX     = 8192
	TEMPERATURE   = 1.0
)

type ChatCompletionRequest struct {
	Model    string                 `json:"model"`
	Messages []*Message             `json:"messages"`
	Params   map[string]interface{} `json:"-"`
}

type CompletionRequest struct {
	Model  string                 `json:"model"`
	Prompt string                 `json:"prompt"`
	Stop   []string               `json:"stop,omitempty"`
	Params map[string]interface{} `json:"-"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
	Logprobs     *any    `json:"logprobs"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResponse struct {
	ID      string    `json:"id"`
	Choices []ChoiceT `json:"choices"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Object  string    `json:"object"`
	Usage   Usage     `json:"usage"`
}

type ChoiceT struct {
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	Text         string `json:"text"`
	Logprobs     *any   `json:"logprobs"`
}

func (r ChatCompletionRequest) MarshalJSON() ([]byte, error) {
	type Alias ChatCompletionRequest
	aux := struct {
		Alias
	}{
		Alias: (Alias)(r),
	}

	data, err := json.Marshal(aux.Alias)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	for k, v := range r.Params {
		result[k] = v
	}

	return json.Marshal(result)
}

func (r CompletionRequest) MarshalJSON() ([]byte, error) {
	type Alias CompletionRequest
	aux := struct {
		Alias
	}{
		Alias: (Alias)(r),
	}

	data, err := json.Marshal(aux.Alias)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	for k, v := range r.Params {
		result[k] = v
	}

	return json.Marshal(result)
}

func CreateChatCompletion(m []*Message) (ChatCompletionResponse, error) {
	var key string
	var requestData ChatCompletionRequest
	var err error
	var jsonData []byte
	var req *http.Request
	var client *http.Client
	var resp *http.Response
	var response ChatCompletionResponse

	key = os.Getenv("API_KEY")
	requestData = ChatCompletionRequest{
		Model:    MODEL,
		Messages: m,
		Params: map[string]interface{}{
			"temperature": TEMPERATURE,
		},
	}

	jsonData, err = json.Marshal(requestData)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return response, err
	}

	req, err = http.NewRequest("POST", CHAT_ENDPOINT, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return response, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return response, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("Error decoding response: %v", err)
		return response, err
	}

	return response, nil
}

func CreateCompletion(prompt string, stop []string) (CompletionResponse, error) {
	var key string
	var requestData CompletionRequest
	var err error
	var jsonData []byte
	var req *http.Request
	var client *http.Client
	var resp *http.Response
	var response CompletionResponse

	key = os.Getenv("API_KEY")
	requestData = CompletionRequest{
		Model:  MODEL,
		Prompt: prompt,
		Stop:   stop,
		Params: map[string]interface{}{
			"temperature": TEMPERATURE,
		},
	}

	jsonData, err = json.Marshal(requestData)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return response, err
	}

	req, err = http.NewRequest("POST", COMP_ENDPOINT, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return response, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return response, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("Error decoding response: %v", err)
		return response, err
	}

	return response, nil
}
