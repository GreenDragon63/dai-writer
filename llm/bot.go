package llm

import (
	"dai-writer/auth"
	"dai-writer/models"
	"log"
	"os"
	"strings"
)

const (
	MODEL_CTX = 2048
	RESPONSE  = 300
)

func replacePlaceholders(s, char, user string) string {
	var replaced string
	replaced = strings.ReplaceAll(s, "{{user}}", user)
	replaced = strings.ReplaceAll(replaced, "You's", "Your")
	replaced = strings.ReplaceAll(replaced, "You is", "You are")
	replaced = strings.ReplaceAll(replaced, "{{char}}", char)
	return replaced
}

func cleanOutput(s, name string) string {
	var replaced string

	replaced = strings.ReplaceAll(s, "\r", "")
	replaced = strings.ReplaceAll(replaced, "<START>", "")
	replaced = strings.ReplaceAll(replaced, "<END>", "")
	replaced = strings.ReplaceAll(replaced, name+":", "")
	replaced = strings.ReplaceAll(replaced, "\\", "")
	replaced = strings.TrimLeft(replaced, "\n")
	replaced = strings.TrimRight(replaced, "\n")
	return replaced
}

func formatContent(prefix, content, separator string) string {
	if len(content) == 0 {
		return ""
	}
	return prefix + content + separator
}

func Generate(u *auth.User, bookId, sceneId, characterId, lineId int) string {
	var stopStrings, words []string
	var debug, stopString, user, memory, newText, streamedText string
	var chara, ch *models.Character
	var scene *models.Scene
	var memorySize, freeSize, responseSize, cid int
	var finished, ok bool
	var err error

	debug = os.Getenv("DEBUG")
	stopStrings, err = loadStopStrings()
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	if debug == "true" {
		log.Println("Stop Strings:")
		for _, stopString = range stopStrings {
			log.Println("|" + stopString + "|")
		}
	}
	finished = false
	newText = ""
	chara, ok = models.LoadCharacter(u, characterId)
	if ok != true {
		log.Printf("Cannot find character %d/%d\n", u.Id, characterId)
		return ""
	}
	chara.Name = strings.Split(chara.Name, "|")[0]
	scene, ok = models.LoadScene(u, bookId, sceneId)
	if ok != true {
		log.Printf("Cannot find scene %d/%d/%d\n", u.Id, bookId, sceneId)
		return ""
	}
	if len(scene.Lines) == 1 {
		user = "You"
		for _, cid = range scene.Characters {
			ch, ok = models.LoadCharacter(u, cid)
			if ok != true {
				log.Printf("Cannot find character %d/%d\n", u.Id, cid)
				return ""
			}
			if ch.IsHuman {
				user = strings.Split(ch.Name, "|")[0]
				break
			}
		}
		return replacePlaceholders(chara.FirstMes, chara.Name, user)
	}

	responseSize = RESPONSE
	for finished == false {
		memorySize = MODEL_CTX - responseSize
		freeSize = responseSize
		memory = botMemory(u, bookId, sceneId, characterId, lineId, memorySize)
		if debug == "true" {
			log.Println(memory)
		}
		for freeSize > 0 {
			if newText != "" {
				words = strings.Split(newText, " ")
				newText = strings.Join(words[:len(words)-1], " ")
			}
			streamedText, finished = GetCompletion(memory + newText)
			newText += streamedText
			freeSize -= 50
			for _, stopString = range stopStrings {
				if strings.Contains(newText, stopString) {
					newText = strings.Split(newText, stopString)[0]
					finished = true
					break
				}
			}
			if debug == "true" {
				log.Println(newText)
			}
			if finished == true {
				break
			}
		}
		responseSize += 100
	}
	return cleanOutput(newText, chara.Name)
}

func botMemory(u *auth.User, bookId, sceneId, characterId, lineId, size int) string {
	var debug, model, user, ltm, stm, currentLine, chatSeparator, exampleSeparator string
	var cid, ltmLength, stmLength, currentLength, lineLength, chatSeparatorLength, exampleSeparatorLength, lastLength int
	var chara *models.Character
	var scene *models.Scene
	var line *models.Line
	var characters map[int]*models.Character
	var ok, insideStm bool
	var err error
	var promptConfig *Prompt

	characters = make(map[int]*models.Character)
	debug = os.Getenv("DEBUG")
	model = GetModel()
	promptConfig, err = loadPromptFormat(model)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	if debug == "true" {
		log.Println("Prompt Config: " + promptConfig.Name)
	}
	chara, ok = models.LoadCharacter(u, characterId)
	if ok != true {
		log.Printf("Cannot find character %d/%d\n", u.Id, characterId)
		return ""
	}
	chara.Name = strings.Split(chara.Name, "|")[0]
	scene, ok = models.LoadScene(u, bookId, sceneId)
	if ok != true {
		log.Printf("Cannot find scene %d/%d/%d\n", u.Id, bookId, sceneId)
		return ""
	}
	user = "You"
	for _, cid = range scene.Characters {
		characters[cid], ok = models.LoadCharacter(u, cid)
		if ok != true {
			log.Printf("Cannot find character %d/%d\n", u.Id, cid)
			return ""
		}
		if characters[cid].IsHuman {
			user = strings.Split(characters[cid].Name, "|")[0]
			break
		}
	}

	ltm = promptConfig.SystemInputSequence + promptConfig.SystemPrompt + promptConfig.SystemOutputSequence
	ltm += formatContent(promptConfig.Description, chara.Description, "\n")
	ltm += formatContent(promptConfig.Personality, chara.Personality, "\n")
	ltm += formatContent(promptConfig.Scenario, chara.Scenario, "\n")
	ltm = replacePlaceholders(ltm, chara.Name, user)
	chatSeparator = replacePlaceholders(promptConfig.ChatSeparator, chara.Name, user)
	exampleSeparator = replacePlaceholders(promptConfig.ExampleSeparator, chara.Name, user)

	ltmLength = GetTokens(ltm)
	chatSeparatorLength = GetTokens(chatSeparator)
	exampleSeparatorLength = GetTokens(exampleSeparator)
	lastLength = GetTokens(promptConfig.OutputSequence + chara.Name + ": ")
	stmLength = size - (ltmLength + chatSeparatorLength + lastLength)

	currentLength = 0
	stm = ""
	insideStm = false
	for i := len(scene.Lines) - 1; i >= 0; i-- {
		if insideStm == false {
			if scene.Lines[i] == lineId {
				insideStm = true
			}
			continue
		}
		line, ok = models.LoadLine(u, bookId, sceneId, scene.Lines[i])
		if ok != true {
			log.Printf("Cannot find line %d/%d/%d/%d\n", u.Id, bookId, sceneId, scene.Lines[i])
			return ""
		}
		if line.Displayed {
			if line.CharacterId == characterId {
				currentLine = promptConfig.OutputSequence + chara.Name + ": " + line.Content[line.Current] + promptConfig.StopSequence
			} else {
				if characters[line.CharacterId] == nil {
					characters[line.CharacterId], ok = models.LoadCharacter(u, line.CharacterId)
					if ok != true {
						log.Printf("Cannot find character %d/%d\n", u.Id, line.CharacterId)
						return ""
					}
				}
				currentLine = promptConfig.InputSequence + strings.Split(characters[line.CharacterId].Name, "|")[0] + ": " + line.Content[line.Current] + promptConfig.StopSequence
			}
			lineLength = GetTokens(currentLine)
			if (currentLength + lineLength) > stmLength {
				break
			}
			currentLength += lineLength
			stm = currentLine + stm
		}
	}
	lineLength = GetTokens(chara.MesExample)
	if (currentLength + lineLength + exampleSeparatorLength) < stmLength {
		stm = formatContent(exampleSeparator, chara.MesExample, "") + formatContent(chatSeparator, stm, "")
	} else {
		stm = formatContent(chatSeparator, stm, "")
	}
	stm = replacePlaceholders(stm, chara.Name, user)
	return ltm + stm + promptConfig.OutputSequence + chara.Name + ": "
}
