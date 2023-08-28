package llm

import (
	"dai-writer/auth"
	"dai-writer/models"
	"fmt"
	"log"
	"os"
	"strings"
)

var promptConfig *Prompt = nil

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

func formatContent(prefix, content string) string {
	if len(content) == 0 {
		return ""
	}
	return prefix + content + "\n"
}

func Generate(u *auth.User, book_id, scene_id, character_id, line_id int) string {
	var stopStrings, words []string
	var stopString, name, memory, new_text, streamed_text string
	var memory_size, free_size, response_size int
	var finished bool

	stopStrings = []string{"You:", "You :", "user:", "USER:"}
	debug := os.Getenv("DEBUG")
	finished = false
	new_text = ""
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

	if len(scene.Lines) == 1 {
		return replacePlaceholders(chara.First_mes, chara.Name)
	}

	response_size = RESPONSE
	for finished == false {
		memory_size = MODEL_CTX - response_size
		free_size = response_size
		memory = botMemory(u, book_id, scene_id, character_id, line_id, memory_size)
		if debug == "true" {
			log.Println(memory)
		}
		for free_size > 0 {
			if new_text != "" {
				words = strings.Split(new_text, " ")
				new_text = strings.Join(words[:len(words)-1], " ")
			}
			streamed_text, finished = GetCompletion(memory + new_text)
			new_text += streamed_text
			free_size -= 50
			for _, stopString = range stopStrings {
				if strings.Contains(new_text, stopString) {
					new_text = strings.Split(new_text, stopString)[0]
					finished = true
					break
				}
			}
			if debug == "true" {
				log.Println(new_text)
			}
			if finished == true {
				break
			}
		}
		response_size += 100
	}
	return cleanOutput(new_text, chara.Name)
}

func botMemory(u *auth.User, book_id, scene_id, character_id, line_id, size int) string {
	var name, ltm, stm, currentLine, model, chatSeparator, exampleSeparator string
	var ltm_length, stmLength, currentLength, lineLength, chatSepLen, exampleSepLen, lastLen int
	var chara *models.Character
	var scene *models.Scene
	var line *models.Line
	var ok, isPygmalion, insideStm bool
	var err error

	if promptConfig == nil {
		promptConfig, err = loadPrompt(os.Getenv("PROMPT"))
		if err != nil {
			log.Println(err.Error())
			return ""
		}
	}
	chara, ok = models.LoadCharacter(u, character_id)
	if ok != true {
		log.Printf("Cannot find character %d\n", character_id)
		return ""
	}
	name = strings.Split(chara.Name, "|")[0]
	chara.Name = name
	scene, ok = models.LoadScene(u, book_id, scene_id)
	if ok != true {
		log.Printf("Cannot find scene %d\n", scene_id)
		return ""
	}
	model = GetModel()
	model = strings.ToLower(model)
	if strings.Contains(model, "pygmalion") {
		isPygmalion = true
		chatSeparator = "<START>\n"
		exampleSeparator = "<START>\n"
		ltm = formatContent(fmt.Sprintf("%s's Persona: ", chara.Name), chara.Description)
		ltm += formatContent("Personality: ", chara.Personality)
		ltm += formatContent("Scenario: ", chara.Scenario)
		lastLen = GetTokens(chara.Name + ": ")
	} else {
		isPygmalion = false
		chatSeparator = fmt.Sprintf("\nThen the roleplay chat between you and %s begins.\n", chara.Name)
		exampleSeparator = fmt.Sprintf("This is how %s should talk\n", chara.Name)
		ltm = formatContent(promptConfig.SystemPrompt+"\n", chara.Description)
		ltm += formatContent(fmt.Sprintf("%s's personality: ", chara.Name), chara.Personality)
		ltm += formatContent("Circumstances and context of the dialogue: ", chara.Scenario)
		lastLen = GetTokens(promptConfig.OutputSequence + "\n" + chara.Name + ": ")
	}
	chatSepLen = GetTokens(chatSeparator)
	exampleSepLen = GetTokens(exampleSeparator)
	ltm = replacePlaceholders(ltm, chara.Name)
	ltm_length = GetTokens(ltm)

	stmLength = size - (ltm_length + chatSepLen + lastLen)
	currentLength = 0
	stm = ""
	insideStm = false
	for i := len(scene.Lines) - 1; i >= 0; i-- {
		if insideStm == false {
			if scene.Lines[i] == line_id {
				insideStm = true
			}
			continue
		}
		line, ok = models.LoadLine(u, book_id, scene_id, scene.Lines[i])
		if ok != true {
			log.Printf("Cannot find scene %d\n", scene_id)
			return ""
		}
		if line.Displayed {
			if line.CharacterId == character_id {
				if isPygmalion {
					currentLine = ""
				} else {
					currentLine = promptConfig.OutputSequence + "\n"
				}
				currentLine += chara.Name + ": " + line.Content[line.Current]
			} else {
				if isPygmalion {
					currentLine = ""
				} else {
					currentLine = promptConfig.InputSequence + "\n"
				}
				currentLine += "You: " + line.Content[line.Current]
			}
			lineLength = GetTokens(currentLine)
			if (currentLength + lineLength) > stmLength {
				break
			}
			currentLength += lineLength
			stm = currentLine + "\n" + stm
		}
	}
	lineLength = GetTokens(chara.Mes_example)
	if (currentLength + lineLength + exampleSepLen) < stmLength {
		stm = exampleSeparator + chara.Mes_example + chatSeparator + stm
	} else {
		stm = chatSeparator + stm
	}
	stm = replacePlaceholders(stm, chara.Name)
	if isPygmalion {
		return ltm + stm + chara.Name + ": "
	} else {
		return ltm + stm + promptConfig.OutputSequence + "\n" + chara.Name + ": "
	}
}
