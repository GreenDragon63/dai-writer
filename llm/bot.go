package llm

import (
	"dai-writer/auth"
	"dai-writer/models"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	MODEL_CTX = 2048
	RESPONSE  = 300
)

func replacePlaceholders(s, name, user string) string {
	var replaced string
	log.Println("Replacing: " + s)
	replaced = strings.ReplaceAll(s, "{{user}}", user)
	replaced = strings.ReplaceAll(replaced, "You's", "Your")
	replaced = strings.ReplaceAll(replaced, "You is", "You are")
	replaced = strings.ReplaceAll(replaced, "{{char}}", name)
	log.Println("Replaced: " + replaced)
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

func formatContent(prefix, content string) string {
	if len(content) == 0 {
		return ""
	}
	return prefix + content + "\n"
}

func Generate(u *auth.User, bookId, sceneId, characterId, lineId int) string {
	var stopStrings, words []string
	var debug, stopString, name, user, memory, newText, streamedText string
	var chara, ch *models.Character
	var memorySize, freeSize, responseSize, cid int
	var finished bool
	var err error

	debug = os.Getenv("DEBUG")
	stopStrings, err = loadStopStrings(os.Getenv("STOP_STRINGS"))
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
	chara, ok := models.LoadCharacter(u, characterId)
	if ok != true {
		log.Printf("Cannot find character %d/%d\n", u.Id, characterId)
		return ""
	}
	name = strings.Split(chara.Name, "|")[0]
	chara.Name = name
	log.Println("Chara name " + chara.Name)
	scene, ok := models.LoadScene(u, bookId, sceneId)
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
		log.Println(("Chara name " + chara.Name))
		log.Println("First message" + chara.FirstMes)
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
	var debug, name, user, ltm, stm, currentLine, model, chatSeparator, exampleSeparator string
	var cid, ltmLength, stmLength, currentLength, lineLength, chatSepLen, exampleSepLen, lastLen int
	var chara *models.Character
	var scene *models.Scene
	var line *models.Line
	var characters map[int]*models.Character
	var ok, isPygmalion, insideStm bool
	var err error
	var promptConfig *Prompt

	characters = make(map[int]*models.Character)
	debug = os.Getenv("DEBUG")
	promptConfig, err = loadPrompt(os.Getenv("PROMPT"))
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
	name = strings.Split(chara.Name, "|")[0]
	chara.Name = name
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
	ltm = replacePlaceholders(ltm, chara.Name, user)
	ltmLength = GetTokens(ltm)

	stmLength = size - (ltmLength + chatSepLen + lastLen)
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
				if isPygmalion {
					currentLine = ""
				} else {
					currentLine = promptConfig.OutputSequence + "\n"
				}
				currentLine += chara.Name + ": " + line.Content[line.Current]
			} else {
				if isPygmalion {
					currentLine = ""
					currentLine += "You: " + line.Content[line.Current]
				} else {
					currentLine = promptConfig.InputSequence + "\n"
					if characters[line.CharacterId] == nil {
						characters[line.CharacterId], ok = models.LoadCharacter(u, line.CharacterId)
						if ok != true {
							log.Printf("Cannot find character %d/%d\n", u.Id, line.CharacterId)
							return ""
						}
					}
					currentLine += strings.Split(characters[line.CharacterId].Name, "|")[0] + ": " + line.Content[line.Current]
				}

			}
			lineLength = GetTokens(currentLine)
			if (currentLength + lineLength) > stmLength {
				break
			}
			currentLength += lineLength
			stm = currentLine + "\n" + stm
		}
	}
	lineLength = GetTokens(chara.MesExample)
	if (currentLength + lineLength + exampleSepLen) < stmLength {
		stm = exampleSeparator + chara.MesExample + chatSeparator + stm
	} else {
		stm = chatSeparator + stm
	}
	stm = replacePlaceholders(stm, chara.Name, user)
	if isPygmalion {
		return ltm + stm + chara.Name + ": "
	} else {
		return ltm + stm + promptConfig.OutputSequence + "\n" + chara.Name + ": "
	}
}
