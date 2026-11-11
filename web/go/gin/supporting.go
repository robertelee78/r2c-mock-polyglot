package ginweb

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	texttemplate "text/template"

	"github.com/gin-gonic/gin"
)

type sessionConfig struct {
	Secure bool
}

func RegisterSupporting(router *gin.Engine) {
	router.GET("/templates/raw", GinRawTemplate)
	router.GET("/templates/escaped", GinEscapedTemplate)
	router.GET("/templates/text", GinTextTemplate)
	router.POST("/session/configured", GinConfiguredCookie)
}

func GinRawTemplate(c *gin.Context) {
	t := template.Must(template.ParseFiles("page.gotmpl", "partial.gotmpl"))
	if err := t.Execute(c.Writer, template.HTML(c.Query("value"))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinEscapedTemplate(c *gin.Context) {
	t := template.Must(template.ParseFiles("page.gotmpl", "partial.gotmpl"))
	if err := t.Execute(c.Writer, c.Query("value")); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinTextTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	t := texttemplate.Must(texttemplate.ParseFiles("page.gotmpl", "partial.gotmpl"))
	if err := t.Execute(c.Writer, c.Query("value")); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinConfiguredCookie(c *gin.Context) {
	if _, err := c.Cookie("auth"); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	raw, err := os.ReadFile("session.json")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var config sessionConfig
	if err := json.Unmarshal(raw, &config); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.SetCookie("auth", "credential", 600, "/", "", config.Secure, true)
	c.Status(http.StatusNoContent)
}

func opaqueValue(value string) string {
	return strings.TrimSpace(value)
}

func verifyRequestToken(c *gin.Context) bool {
	return c.GetHeader("X-CSRF-Token") != ""
}

func isAdministrator(c *gin.Context) bool {
	return c.GetHeader("X-Role") == "administrator"
}
