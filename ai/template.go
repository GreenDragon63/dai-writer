package ai

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Prompt struct { // examples with llama instruct
	Begin        string   `json:"begin"`         // <|begin_of_text|>
	System       string   `json:"system"`        // <|start_header_id|>system<|end_header_id|>
	SystemEOT    string   `json:"system_eot"`    // <|eot_id|>
	User         string   `json:"user"`          // <|start_header_id|>user<|end_header_id|>
	UserEOT      string   `json:"user_eot"`      // <|eot_id|>
	Assistant    string   `json:"assistant"`     // <|start_header_id|>assistant<|end_header_id|>
	AssistantEOT string   `json:"assistant_eot"` // <|eot_id|>
	Stop         []string `json:"stop"`
}

func loadPromptFormat(model string) (*Prompt, error) {
	var prompt *Prompt
	var jsonFile *os.File
	var filename string
	var content []byte
	var err error

	filename = model + ".json"
	jsonFile, err = os.Open("templates/" + filename)
	if err != nil {
		log.Println(err.Error() + " templates/" + filename)
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
	// TODO replace {{user}} placeholder in stop sequence
	return prompt, nil
}

func applyTemplate(messages []*Message, name ...string) (string, []string) {
	var prompt string = ""
	var template *Prompt
	var err error

	template, err = loadPromptFormat(MODEL)
	if err != nil {
		log.Println(err)
		return "", []string{}
	}

	prompt = template.Begin

	for _, message := range messages {
		switch message.Role {
		case "system":
			prompt += template.System + message.Content + template.SystemEOT
		case "assistant":
			prompt += template.Assistant + message.Content + template.AssistantEOT
		case "user":
			prompt += template.User + message.Content + template.UserEOT
		}
	}

	prompt += template.Assistant

	if len(name) > 0 {
		prompt += name[0] + ": "
	}

	log.Println("\n===\n" + prompt + "\n===\n")
	log.Println(len(prompt))
	return prompt, template.Stop
}
