package echoweb

import (
	"html"
	"html/template"

	"github.com/labstack/echo/v4"
)

func RegisterTypedContexts(e *echo.Echo) {
	e.GET("/js/positive", EchoJsPositive)
	e.GET("/js/control", EchoJsControl)
	e.GET("/js/uncertain", EchoJsUncertain)
	e.GET("/js/html-encoder", EchoJsWrongHtmlEncoder)
	e.GET("/url/positive", EchoUrlPositive)
	e.GET("/url/control", EchoUrlControl)
	e.GET("/url/uncertain", EchoUrlUncertain)
	e.GET("/url/html-encoder", EchoUrlWrongHtmlEncoder)
}

func EchoJsPositive(c echo.Context) error {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	return t.Execute(c.Response(), template.JS(c.QueryParam("value")))
}

func EchoJsControl(c echo.Context) error {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	return t.Execute(c.Response(), template.JS("const answer = 1"))
}

func EchoJsUncertain(c echo.Context) error {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	return t.Execute(c.Response(), template.JS(opaqueValue(c.QueryParam("value"))))
}

func EchoJsWrongHtmlEncoder(c echo.Context) error {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	return t.Execute(c.Response(), template.JS(html.EscapeString(c.QueryParam("value"))))
}

func EchoUrlPositive(c echo.Context) error {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	return t.Execute(c.Response(), template.URL(c.QueryParam("value")))
}

func EchoUrlControl(c echo.Context) error {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	return t.Execute(c.Response(), template.URL("https://trusted.example"))
}

func EchoUrlUncertain(c echo.Context) error {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	return t.Execute(c.Response(), template.URL(opaqueValue(c.QueryParam("value"))))
}

func EchoUrlWrongHtmlEncoder(c echo.Context) error {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	return t.Execute(c.Response(), template.URL(html.EscapeString(c.QueryParam("value"))))
}
