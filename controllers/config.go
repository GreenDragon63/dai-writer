package controllers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	var tmplContent []byte
	var err error
	var buf bytes.Buffer
	var tmpl *template.Template

	c.Header("Content-Type", "application/javascript")

	tmplContent, err = os.ReadFile("js/config.tmpl")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template file not found")
		return
	}

	// Données à passer au modèle
	data := gin.H{
		"prefix": os.Getenv("URL_PREFIX"),
	}

	tmpl, err = template.New("js").Parse(string(tmplContent))
	if err != nil {
		c.String(http.StatusInternalServerError, "Template syntax error")
		return
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execution error")
		return
	}

	c.String(http.StatusOK, buf.String())
}
