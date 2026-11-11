package echoweb

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	texttemplate "text/template"

	"github.com/labstack/echo/v4"
)

type sessionConfig struct {
	Secure bool
}

func RegisterSupporting(e *echo.Echo) {
	e.GET("/templates/raw", EchoRawTemplate)
	e.GET("/templates/escaped", EchoEscapedTemplate)
	e.GET("/templates/text", EchoTextTemplate)
	e.POST("/session/configured", EchoConfiguredCookie)
}

func EchoRawTemplate(c echo.Context) error {
	t := template.Must(template.ParseFiles("page.gotmpl", "partial.gotmpl"))
	return t.Execute(c.Response(), template.HTML(c.QueryParam("value")))
}

func EchoEscapedTemplate(c echo.Context) error {
	t := template.Must(template.ParseFiles("page.gotmpl", "partial.gotmpl"))
	return t.Execute(c.Response(), c.QueryParam("value"))
}

func EchoTextTemplate(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	t := texttemplate.Must(texttemplate.ParseFiles("page.gotmpl", "partial.gotmpl"))
	return t.Execute(c.Response(), c.QueryParam("value"))
}

func EchoConfiguredCookie(c echo.Context) error {
	if _, err := c.Cookie("auth"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	raw, err := os.ReadFile("session.json")
	if err != nil {
		return err
	}
	var config sessionConfig
	if err := json.Unmarshal(raw, &config); err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{Name: "auth", Value: "credential", Secure: config.Secure, HttpOnly: true})
	return c.NoContent(http.StatusNoContent)
}

func opaqueValue(value string) string {
	return strings.TrimSpace(value)
}

func verifyRequestToken(c echo.Context) bool {
	return c.Request().Header.Get("X-CSRF-Token") != ""
}

func isAdministrator(c echo.Context) bool {
	return c.Request().Header.Get("X-Role") == "administrator"
}
