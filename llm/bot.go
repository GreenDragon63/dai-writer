package llm

import (
	"dai-writer/ai"
	"dai-writer/auth"
	"dai-writer/engram"
	"dai-writer/models"

	"log"
	"os"
	"strconv"
	"strings"
)

const (
	STM_SIZE    = 1024
	ENGRAN_SIZE = 1024
	RESPONSE    = 1024
)

func getTokens(s string) int {
	// Fake token count for simplicity
	// Assuming 4 characters per token
	return len(s) / 4
}

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

func Generate(u *auth.User, bookId, sceneId, characterId, lineId int, input string) string {
	var response ai.CompletionResponse
	var stopStrings []string
	var debug, stopString, user, memory, newText string
	var chara, ch *models.Character
	var scene *models.Scene
	var memorySize, cid int
	var ok bool
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

	memorySize = ai.MODEL_CTX - RESPONSE
	memory = botMemory(u, bookId, sceneId, characterId, lineId, memorySize, input)
	response, err = ai.CreateCompletion(memory, stopStrings)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	newText = response.Choices[0].Text

	return cleanOutput(newText, chara.Name)
}

func botMemory(u *auth.User, bookId, sceneId, characterId, lineId, size int, input string) string {
	var debug, model, user, ltm, engrm, stm, currentLine, chatSeparator, exampleSeparator, systemInput string
	var i, bId, sId, lId, cid, ltmLength, engramLength, stmLength, currentLength, lineLength, chatSeparatorLength, exampleSeparatorLength, lastLength, systemInputLength int
	var chara *models.Character
	var scene *models.Scene
	var line *models.Line
	var characters map[int]*models.Character
	var ok, insideStm, engramFull bool
	var err error
	var promptConfig *Prompt
	var answer []map[string]float64

	characters = make(map[int]*models.Character)
	debug = os.Getenv("DEBUG")
	model = ai.MODEL
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
	if len(scene.Description) > 0 {
		ltm += formatContent(promptConfig.Scenario, scene.Description, "\n")
	} else {
		ltm += formatContent(promptConfig.Scenario, chara.Scenario, "\n")
	}
	ltm = replacePlaceholders(ltm, chara.Name, user)
	chatSeparator = replacePlaceholders(promptConfig.ChatSeparator, chara.Name, user)
	exampleSeparator = replacePlaceholders(promptConfig.ExampleSeparator, chara.Name, user)

	systemInput = formatContent(promptConfig.SystemInputSequence, input, promptConfig.SystemOutputSequence)
	if len(systemInput) > 0 {
		systemInputLength = getTokens(systemInput)
	} else {
		systemInputLength = 0
	}

	ltmLength = getTokens(ltm)
	chatSeparatorLength = getTokens(chatSeparator)
	exampleSeparatorLength = getTokens(exampleSeparator)
	lastLength = getTokens(promptConfig.OutputSequence + chara.Name + ": ")
	stmLength = size - (ltmLength + chatSeparatorLength + lastLength + systemInputLength) // TODO refactor
	engramLength = ENGRAN_SIZE
	stmLength = STM_SIZE

	// Create engram
	currentLength = 0
	engrm = ""
	insideStm = false
	engramFull = false
	for i = len(scene.Lines) - 1; i >= 0; i-- {
		if insideStm == false {
			if scene.Lines[i] == lineId {
				insideStm = true
			}
			continue
		}
		answer, ok = engram.Search(u, bookId, sceneId, scene.Lines[i], chara.Id)
		for _, match := range answer {
			for ubsl, score := range match {
				if score > 0.7 { // TODO: make this configurable
					ubslDecoded := strings.Split(ubsl, "-")
					if len(ubslDecoded) != 4 {
						log.Printf("Invalid UBSL: %s\n", ubsl)
						return ""
					}
					bId, err = strconv.Atoi(ubslDecoded[1])
					if err != nil {
						log.Printf("Invalid book id: %s\n", ubslDecoded[1])
						return ""
					}
					sId, err = strconv.Atoi(ubslDecoded[2])
					if err != nil {
						log.Printf("Invalid scene id: %s\n", ubslDecoded[2])
						return ""
					}
					lId, err = strconv.Atoi(ubslDecoded[3])
					if err != nil {
						log.Printf("Invalid line id: %s\n", ubslDecoded[3])
						return ""
					}
					line, ok = models.LoadLine(u, bId, sId, lId)
					if ok != true {
						log.Printf("Cannot find line %d/%d/%d/%d\n", u.Id, bookId, sceneId, scene.Lines[i])
						return ""
					}
					if line.Displayed {
						if line.CharacterId == characterId {
							currentLine = chara.Name + ": " + line.Content[line.Current] + promptConfig.StopSequence
						} else {
							if characters[line.CharacterId] == nil {
								characters[line.CharacterId], ok = models.LoadCharacter(u, line.CharacterId)
								if ok != true {
									log.Printf("Cannot find character %d/%d\n", u.Id, line.CharacterId)
									return ""
								}
							}
							currentLine = strings.Split(characters[line.CharacterId].Name, "|")[0] + ": " + line.Content[line.Current] + promptConfig.StopSequence
						}
						lineLength = getTokens(currentLine)
						if (currentLength + lineLength) > engramLength {
							engramFull = true
							break
						}
						if strings.Contains(engrm, currentLine) {
							continue
						}
						currentLength += lineLength
						engrm = currentLine + engrm
					}
				}
			}
		}
		if engramFull == true {
			break
		}
	}
	log.Println(engrm)

	// Create STM
	currentLength = 0
	stm = ""
	insideStm = false
	for i = len(scene.Lines) - 1; i >= 0; i-- {
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
			lineLength = getTokens(currentLine)
			if (currentLength + lineLength) > stmLength {
				break
			}
			currentLength += lineLength
			stm = currentLine + stm
		}
	}
	lineLength = getTokens(chara.MesExample)
	if (currentLength + lineLength + exampleSeparatorLength) < stmLength {
		stm = formatContent(exampleSeparator, chara.MesExample, "") + formatContent(chatSeparator, stm, "")
	} else {
		stm = formatContent(chatSeparator, stm, "")
	}
	stm = replacePlaceholders(stm, chara.Name, user)
	return ltm + stm + promptConfig.OutputSequence + systemInput + chara.Name + ": "
}
