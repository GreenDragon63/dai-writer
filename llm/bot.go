package llm

import (
	"dai-writer/auth"
	"dai-writer/models"
	"fmt"
	"log"
	"strings"
)

const (
	MODEL_CTX = 2048
	RESPONSE  = 300
)

func replacePlaceholders(s, name string) string {
	replaced := strings.ReplaceAll(s, "{{user}}", "You")
	replaced = strings.ReplaceAll(replaced, "You's", "Your")
	replaced = strings.ReplaceAll(replaced, "You is", "You are")
	replaced = strings.ReplaceAll(replaced, "{{char}}", name)
	return replaced
}

func cleanOutput(s, name string) string {
	replaced := strings.ReplaceAll(s, "\r", "")
	replaced = strings.ReplaceAll(replaced, "<START>", "")
	replaced = strings.ReplaceAll(replaced, "<END>", "")
	replaced = strings.ReplaceAll(replaced, name+":", "")
	replaced = strings.ReplaceAll(replaced, "\\", "")
	replaced = strings.TrimLeft(replaced, "\n")
	replaced = strings.TrimRight(replaced, "\n")
	return replaced
}

func Generate(u *auth.User, book_id, scene_id, character_id, line_id int) string {
	var words []string
	var memory, new_text, streamed_text string
	var memory_size, free_size, response_size int
	var finished bool

	finished = false
	new_text = ""
	chara, ok := models.LoadCharacter(u, character_id)
	if ok != true {
		log.Printf("Cannot find character %d\n", character_id)
		return ""
	}
	scene, ok := models.LoadScene(u, book_id, scene_id)
	if ok != true {
		log.Printf("Cannot find scene %d\n", scene_id)
		return ""
	}

	if len(scene.Lines) == 1 {
		return cleanOutput(chara.First_mes, chara.Name)
	}

	response_size = RESPONSE
	for finished == false {
		memory_size = MODEL_CTX - response_size
		free_size = response_size
		memory = botMemory(u, book_id, scene_id, character_id, line_id, memory_size)
		log.Println(memory)
		for free_size > 0 {
			if new_text != "" {
				words = strings.Split(new_text, " ")
				new_text = strings.Join(words[:len(words)-1], " ")
			}
			streamed_text, finished = GetCompletion(memory + new_text)
			new_text += streamed_text
			free_size -= 50
			if strings.Contains(new_text, "You:") {
				new_text = strings.Split(new_text, "You:")[0]
				finished = true
				break
			}
			if strings.Contains(new_text, "You :") {
				new_text = strings.Split(new_text, "You :")[0]
				finished = true
				break
			}
			log.Println(new_text)
			if finished == true {
				break
			}
		}
		response_size += 100
	}
	return cleanOutput(new_text, chara.Name)
}

func botMemory(u *auth.User, book_id, scene_id, character_id, line_id, size int) string {
	var name, ltm, stm, current_line string
	var ltm_length, stm_length, current_length, line_length int

	chara, ok := models.LoadCharacter(u, character_id)
	if ok != true {
		log.Printf("Cannot find character %d\n", character_id)
		return ""
	}
	name = strings.Split(chara.Name, "|")[0]
	chara.Name = name
	scene, ok := models.LoadScene(u, book_id, scene_id)
	if ok != true {
		log.Printf("Cannot find scene %d\n", scene_id)
		return ""
	}
	if chara.Personality != "" {
		ltm = fmt.Sprintf("%s's Persona: %s\nPersonality: %s\nScenario: %s\n", chara.Name, chara.Description, chara.Personality, chara.Scenario)
	} else {
		ltm = fmt.Sprintf("%s's Persona: %s\nScenario: %s\n", chara.Name, chara.Description, chara.Scenario)
	}
	ltm = replacePlaceholders(ltm, chara.Name)
	ltm_length = GetTokens(ltm)
	stm_length = size - ltm_length
	current_length = 0
	stm = ""
	for i := len(scene.Lines) - 1; i >= 0; i-- {
		line, ok := models.LoadLine(u, book_id, scene_id, scene.Lines[i])
		if ok != true {
			log.Printf("Cannot find scene %d\n", scene_id)
			return ""
		}
		if line.Displayed && line.Id != line_id {
			if line.CharacterId == character_id {
				current_line = chara.Name + ": " + line.Content[line.Current]
			} else {
				current_line = "You: " + line.Content[line.Current]
			}
			line_length = GetTokens(current_line)
			if (current_length + line_length) > stm_length {
				break
			}
			current_length += line_length
			stm = current_line + "\n" + stm
		}
	}
	line_length = GetTokens(chara.Mes_example)
	if (current_length + line_length) < stm_length {
		stm = chara.Mes_example + "\n" + stm
	} else {
		stm = "<START>\n" + stm
	}
	stm = replacePlaceholders(stm, chara.Name)
	return ltm + stm + chara.Name + ": "
}
