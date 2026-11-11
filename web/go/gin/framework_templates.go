package ginweb

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FrameworkTemplateRouter() *gin.Engine {
	router := gin.New()
	router.LoadHTMLFiles("framework.html")
	router.GET("/framework/positive", GinFrameworkTemplatePositive)
	router.GET("/framework/control", GinFrameworkTemplateControl)
	router.GET("/framework/uncertain", GinFrameworkTemplateUncertain)
	return router
}

func GinFrameworkTemplatePositive(c *gin.Context) {
	c.HTML(http.StatusOK, "framework.html", template.HTML(c.Query("value")))
}

func GinFrameworkTemplateControl(c *gin.Context) {
	c.HTML(http.StatusOK, "framework.html", c.Query("value"))
}

func GinFrameworkTemplateUncertain(c *gin.Context) {
	c.HTML(http.StatusOK, opaqueValue(c.Query("template")), c.Query("value"))
}
