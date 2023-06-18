package controllers

import (
	"dai-writer/auth"
	"dai-writer/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Book(c *gin.Context) {
	c.HTML(http.StatusOK, "book.tmpl", gin.H{
		"title": "Book",
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
