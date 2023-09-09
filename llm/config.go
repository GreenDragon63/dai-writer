package llm

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

type AutoConfig struct {
	Expression string `json:"expression"`
	Prompt     string `json:"prompt"`
	Model      string `json:"model"`
}

type Model struct {
	Prompt       string  `json:"prompt"`
	MaxNewTokens int     `json:"max_new_tokens"`
	Temperature  float64 `json:"temperature"`
	TopP         float64 `json:"top_p"`
	TopK         int     `json:"top_k"`
	RepPen       float64 `json:"repetition_penalty"`
}

type Prompt struct {
	Name                 string `json:"name"`
	SystemInputSequence  string `json:"system_input_sequence"`
	SystemPrompt         string `json:"system_prompt"`
	SystemOutputSequence string `json:"system_output_sequence"`
	Description          string `json:"description"`
	Personality          string `json:"personality"`
	Scenario             string `json:"scenario"`
	ExampleSeparator     string `json:"example_separator"`
	ChatSeparator        string `json:"chat_separator"`
	InputSequence        string `json:"input_sequence"`
	OutputSequence       string `json:"output_sequence"`
	StopSequence         string `json:"stop_sequence"`
}

func loadModelConfig(model string) (*Model, error) {
	var autoConfig []AutoConfig
	var modelConfig *Model
	var confFile, jsonFile *os.File
	var filename string
	var content []byte
	var err error

	filename = ""
	model = strings.ToLower(model)
	confFile, err = os.Open("config/auto_config.json")
	if err != nil {
		log.Println(err.Error())
		return modelConfig, err
	}
	defer confFile.Close()
	content, err = io.ReadAll(confFile)
	if err != nil {
		log.Println(err.Error())
		return modelConfig, err
	}
	err = json.Unmarshal(content, &autoConfig)
	if err != nil {
		log.Println(err.Error())
		return modelConfig, err
	}
	for _, conf := range autoConfig {
		if strings.Contains(model, conf.Expression) {
			filename = conf.Model
			break
		}
	}
	if filename == "" {
		filename = os.Getenv("DEFAULT_MODEL_PARAMETERS")
	}
	jsonFile, err = os.Open("config/" + filename)
	if err != nil {
		log.Println(err.Error())
		return modelConfig, err
	}
	defer jsonFile.Close()
	content, err = io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err.Error())
		return modelConfig, err
	}
	err = json.Unmarshal(content, &modelConfig)
	if err != nil {
		log.Println(err.Error())
		return modelConfig, err
	}
	return modelConfig, nil
}

func loadPromptFormat(model string) (*Prompt, error) {
	var autoConfig []AutoConfig
	var prompt *Prompt
	var confFile, jsonFile *os.File
	var filename string
	var content []byte
	var err error

	filename = ""
	model = strings.ToLower(model)
	confFile, err = os.Open("config/auto_config.json")
	if err != nil {
		log.Println(err.Error())
		return prompt, err
	}
	defer confFile.Close()
	content, err = io.ReadAll(confFile)
	if err != nil {
		log.Println(err.Error())
		return prompt, err
	}
	err = json.Unmarshal(content, &autoConfig)
	if err != nil {
		log.Println(err.Error())
		return prompt, err
	}
	for _, conf := range autoConfig {
		if strings.Contains(model, conf.Expression) {
			filename = conf.Prompt
			break
		}
	}
	if filename == "" {
		filename = os.Getenv("DEFAULT_PROMPT_FORMAT")
	}
	jsonFile, err = os.Open("config/" + filename)
	if err != nil {
		log.Println(err.Error())
		return prompt, err
	}
	defer jsonFile.Close()
	content, err = io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err.Error())
		return prompt, err
	}
	err = json.Unmarshal(content, &prompt)
	if err != nil {
		log.Println(err.Error())
		return prompt, err
	}
	return prompt, nil
}

func loadStopStrings() ([]string, error) {
	var stopStrings []string
	var file *os.File
	var err error
	var scanner *bufio.Scanner

	file, err = os.Open("config/stop_strings.txt")
	if err != nil {
		log.Println(err.Error())
		return stopStrings, err
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		stopStrings = append(stopStrings, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
		return stopStrings, err
	}
	return stopStrings, nil
}
