package llm

import (
	"dai-writer/auth"
	"dai-writer/models"
	"fmt"
	"log"
	"strings"
)

func Generate(u *auth.User, book_id, scene_id, character_id, line_id int) string {
	var memory string
	var size int

	size = 2048 - 300
	memory = botMemory(u, book_id, scene_id, character_id, line_id, size)
	return memory
}

func replacePlaceholders(s, name string) string {
	replaced := strings.ReplaceAll(s, "{{user}}", "You")
	replaced = strings.ReplaceAll(replaced, "You's", "Your")
	replaced = strings.ReplaceAll(replaced, "You is", "You are")
	replaced = strings.ReplaceAll(replaced, "{{char}}", name)
	return replaced
}

func botMemory(u *auth.User, book_id, scene_id, character_id, line_id, size int) string {
	var ltm, stm, current_line string
	var ltm_length, stm_length, current_length, line_length int

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
			current_line = line.Content[line.Current]
			line_length = GetTokens(current_line)
			if (current_length + line_length) > stm_length {
				break
			}
			current_length += line_length // TODO handle character name
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
	return ltm + stm
}
