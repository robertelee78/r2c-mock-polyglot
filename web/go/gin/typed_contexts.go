package ginweb

import (
	"html"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterTypedContexts(router *gin.Engine) {
	router.GET("/js/positive", GinJsPositive)
	router.GET("/js/control", GinJsControl)
	router.GET("/js/uncertain", GinJsUncertain)
	router.GET("/js/html-encoder", GinJsWrongHtmlEncoder)
	router.GET("/url/positive", GinUrlPositive)
	router.GET("/url/control", GinUrlControl)
	router.GET("/url/uncertain", GinUrlUncertain)
	router.GET("/url/html-encoder", GinUrlWrongHtmlEncoder)
}

func GinJsPositive(c *gin.Context) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(c.Writer, template.JS(c.Query("value"))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinJsControl(c *gin.Context) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(c.Writer, template.JS("const answer = 1")); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinJsUncertain(c *gin.Context) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(c.Writer, template.JS(opaqueValue(c.Query("value")))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinJsWrongHtmlEncoder(c *gin.Context) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(c.Writer, template.JS(html.EscapeString(c.Query("value")))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinUrlPositive(c *gin.Context) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(c.Writer, template.URL(c.Query("value"))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinUrlControl(c *gin.Context) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(c.Writer, template.URL("https://trusted.example")); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinUrlUncertain(c *gin.Context) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(c.Writer, template.URL(opaqueValue(c.Query("value")))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinUrlWrongHtmlEncoder(c *gin.Context) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(c.Writer, template.URL(html.EscapeString(c.Query("value")))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
