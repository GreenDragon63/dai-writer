package auth

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

func Login(c *gin.Context) {
	var userReq UserRequest
	var u User
	if err := c.BindJSON(&userReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":3306)/"+os.Getenv("DB_NAME"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	stmt, err := db.Prepare("SELECT id, username, password, salt FROM users where username=?")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userReq.Username)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	var count = 0
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.Salt)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		passwordOk := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userReq.Password+u.Salt))
		if passwordOk != nil {
			log.Println(passwordOk.Error())
		} else {
			count++
		}
	}
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	randomString, err := GenerateRandomString(64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Random error"})
		return
	}
	c.SetCookie("session", randomString, 60*60*24*365, "", "", false, true)
	_, err = db.Exec("UPDATE users SET session = ? WHERE id = ?", randomString, u.Id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	c.JSON(http.StatusOK, u.Username)
}

func GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("LOCAL_INSTALL") == "true" {
			var u User
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
		if cookie, err := c.Cookie("session"); err == nil {
			db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":3306)/"+os.Getenv("DB_NAME"))
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
				return
			}
			defer db.Close()

			err = db.Ping()
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
				return
			}

			stmt, err := db.Prepare("SELECT id, username, password, salt, level, credits FROM users where session = ?")
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
				return
			}
			defer stmt.Close()

			rows, err := stmt.Query(cookie)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
				return
			}

			for rows.Next() {
				var u User
				err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.Salt, &u.Level, &u.Credits)
				if err != nil {
					log.Println(err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
					return
				}
				c.Set("current_user", u)
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden with no session"})
		c.Abort()
	}
}
