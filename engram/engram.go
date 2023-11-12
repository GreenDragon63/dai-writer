package engram

import (
	"dai-writer/auth"
	"dai-writer/models"
)

func ExportScene(u *auth.User, bookId, sceneId int) bool {
	var lines []*models.Line
	var line *models.Line
	var ok bool

	lines, ok = models.ListLine(u, bookId, sceneId)
	if ok != true {
		return false
	}

	for _, line = range lines {
		if line.Displayed == true {
			ok = ExportLine(u, line)
		} else {
			ok = RemoveLine(u, line)
		}
	}
	return true
}
