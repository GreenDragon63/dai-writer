package llm

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type Prompt struct {
	InputSequence     string `json:"input_sequence"`
	Macro             bool   `json:"macro"`
	Name              string `json:"name"`
	Names             bool   `json:"names"`
	OutputSequence    string `json:"output_sequence"`
	SeparatorSequence string `json:"separator_sequence"`
	StopSequence      string `json:"stop_sequence"`
	SystemPrompt      string `json:"system_prompt"`
	SystemSequence    string `json:"system_sequence"`
	Wrap              bool   `json:"wrap"`
}

func loadPrompt(filename string) (*Prompt, error) {
	var prompt *Prompt
	var jsonFile *os.File
	var content []byte
	var err error

	jsonFile, err = os.Open(filename)
	if err != nil {
		return prompt, err
	}
	defer jsonFile.Close()
	content, err = io.ReadAll(jsonFile)
	if err != nil {
		return prompt, err
	}
	err = json.Unmarshal(content, &prompt)
	if err != nil {
		return prompt, err
	}
	return prompt, nil
}

func loadStopStrings(filename string) ([]string, error) {
	var stopStrings []string
	var file *os.File
	var err error
	var scanner *bufio.Scanner

	file, err = os.Open(filename)
	if err != nil {
		return stopStrings, err
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		stopStrings = append(stopStrings, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return stopStrings, err
	}
	return stopStrings, nil
}