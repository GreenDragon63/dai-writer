package auth

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Level    int    `json:"level"`
	Credits  int    `json:"credits"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const prefixUser string = "Users/"
const prefixSession string = "Sessions/"

func CheckUsername(username string) bool {
	var regex *regexp.Regexp

	regex = regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !regex.MatchString(username) {
		return false
	}
	return true
}

func Login(c *gin.Context) {
	var userReq UserRequest
	var u User
	var filename, randomString string
	var err, passwordOk error
	var content []byte

	if err = c.BindJSON(&userReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	if !CheckUsername(userReq.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		return
	}
	filename = prefixUser + userReq.Username + ".json"
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	content, err = os.ReadFile(filename)
	err = json.Unmarshal(content, &u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	passwordOk = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userReq.Password+u.Salt))
	if passwordOk != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	randomString, err = GenerateRandomString(64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Random error"})
		return
	}
	c.SetCookie("session", randomString, 60*60*24*365, "", "", false, true)
	filename = prefixSession + randomString + ".txt"
	err = os.WriteFile(filename, []byte(u.Username), 0644)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	c.JSON(http.StatusOK, u.Username)
}

func GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u User
		var filename, cookie string
		var username, content []byte
		var err error

		if os.Getenv("LOCAL_INSTALL") == "true" {
			u.Id = 0
			u.Username = "admin"
			u.Password = ""
			u.Salt = ""
			u.Level = 1
			u.Credits = -1
			c.Set("current_user", u)
			c.Next()
			return
		}

		if cookie, err = c.Cookie("session"); err == nil {
			filename = prefixSession + cookie + ".txt"
			username, err = os.ReadFile(filename)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
				return
			}
			filename = prefixUser + string(username) + ".json"
			if _, err = os.Stat(filename); os.IsNotExist(err) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
				return
			}
			content, err = os.ReadFile(filename)
			err = json.Unmarshal(content, &u)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
				return
			}
			c.Set("current_user", u)
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden with no session"})
		c.Abort()
	}
}

func InitUser() {
	var err error

	if os.Getenv("LOCAL_INSTALL") == "true" {
		return
	}

	_, err = os.Stat(prefixUser)
	if os.IsNotExist(err) {
		err = os.Mkdir(prefixUser, 0755)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		err = os.Mkdir(prefixSession, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}
		CreateUser()
	}
}

func CreateUser() {
	var user User
	var id int
	var username, salt, filename string
	var password, saltedPassword, hash []byte
	var files []fs.DirEntry
	var file fs.DirEntry
	var content []byte
	var err error

	fmt.Print("User : ")
	fmt.Scanln(&username)
	fmt.Print("Password : ")
	password, _ = term.ReadPassword(int(os.Stdin.Fd()))

	salt, err = GenerateRandomString(16)
	if err != nil {
		log.Fatal(err.Error())
	}
	saltedPassword = append(password, []byte(salt)...)
	hash, err = bcrypt.GenerateFromPassword(saltedPassword, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err.Error())
	}

	id = 0
	files, err = os.ReadDir(prefixUser)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file = range files {
		content, err = os.ReadFile(prefixUser + file.Name())
		if err != nil {
			log.Fatal(err.Error())
		}
		err = json.Unmarshal(content, &user)
		if err != nil {
			log.Fatal(err.Error())
		}
		if user.Id > id {
			id = user.Id
		}
	}
	id += 1

	user.Id = id
	user.Username = username
	user.Password = string(hash)
	user.Salt = string(salt)
	user.Level = 1
	user.Credits = -1

	filename = prefixUser + username + ".json"
	content, err = json.Marshal(user)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
}
