package controllers

import (
	"dai-writer/auth"
	"dai-writer/models"

	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Book(c *gin.Context) {
	c.HTML(http.StatusOK, "book.tmpl", gin.H{
		"title":  "Book",
		"prefix": os.Getenv("URL_PREFIX"),
		"js":     "book.js",
	})
}

func ListBook(c *gin.Context) {
	var user auth.User

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	user = u.(auth.User)
	book, ok := models.ListBook(&user)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func GetBook(c *gin.Context) {
	var user auth.User
	var id int

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	book, ok := models.LoadBook(&user, id)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func PostBook(c *gin.Context) {
	var user auth.User
	var id int
	var ok bool
	var Book models.Book

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	if err := c.BindJSON(&Book); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book format"})
		return
	}
	ok = models.SaveBook(&user, id, Book)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func DeleteBook(c *gin.Context) {
	var user auth.User
	var id int
	var ok bool

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	ok = models.DeleteBook(&user, id)
	if ok != true {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ExportBook(c *gin.Context) {
	var user auth.User
	var name, content string
	var id int
	var characters map[int]*models.Character

	characters = make(map[int]*models.Character)
	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	book, ok := models.LoadBook(&user, id)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	content = ""
	for s := 0; s < len(book.Scenes); s++ {
		scene, ok := models.LoadScene(&user, id, book.Scenes[s])
		if ok != true {
			log.Printf("Cannot find scene %d\n", book.Scenes[s])
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
		if scene.Displayed {
			for l := 0; l < len(scene.Lines); l++ {
				line, ok := models.LoadLine(&user, id, book.Scenes[s], scene.Lines[l])
				if ok != true {
					log.Printf("Cannot find line %d\n", scene.Lines[l])
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
					return
				}
				if line.Displayed {
					if characters[line.CharacterId] == nil {
						characters[line.CharacterId], ok = models.LoadCharacter(&user, line.CharacterId)
						if ok != true {
							log.Printf("Cannot find character %d\n", line.CharacterId)
							c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
							return
						}
						name = strings.Split(characters[line.CharacterId].Name, "|")[0]
						characters[line.CharacterId].Name = name
					}
					content += characters[line.CharacterId].Name + ": " + line.Content[line.Current] + "\n"
				}
			}
			content += "\n* * *\n\n"
		}
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=export_book_"+c.Param("id")+".txt")
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusOK, content)
}
