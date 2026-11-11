package ginweb

import (
	"html"
	"html/template"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Entity struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Role string `json:"role" form:"role"`
}

var database *gorm.DB

var policy *casbin.Enforcer

func RegisterCases(router *gin.Engine) {
	router.GET("/html/positive", GinHtmlPositive)
	router.GET("/html/control", GinHtmlControl)
	router.GET("/html/uncertain", GinHtmlUncertain)
	router.POST("/csrf/positive", GinCsrfPositive)
	router.POST("/csrf/control", GinCsrfControl)
	router.POST("/csrf/uncertain", GinCsrfUncertain)
	router.GET("/redirect/positive", GinRedirectPositive)
	router.GET("/redirect/control", GinRedirectControl)
	router.GET("/redirect/uncertain", GinRedirectUncertain)
	router.POST("/policy/positive", GinPolicyPositive)
	router.POST("/policy/control", GinPolicyControl)
	router.POST("/policy/uncertain", GinPolicyUncertain)
	router.POST("/session/positive", GinSessionPositive)
	router.POST("/session/control", GinSessionControl)
	router.POST("/session/uncertain", GinSessionUncertain)
	router.GET("/cors/positive", GinCorsPositive)
	router.GET("/cors/control", GinCorsControl)
	router.GET("/cors/uncertain", GinCorsUncertain)
	router.POST("/binding/positive", GinBindingPositive)
	router.POST("/binding/control", GinBindingControl)
	router.POST("/binding/uncertain", GinBindingUncertain)
}

func GinHtmlPositive(c *gin.Context) {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	if err := t.Execute(c.Writer, template.HTML(c.Query("value"))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinHtmlControl(c *gin.Context) {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	if err := t.Execute(c.Writer, template.HTML(html.EscapeString(c.Query("value")))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinHtmlUncertain(c *gin.Context) {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	if err := t.Execute(c.Writer, template.HTML(opaqueValue(c.Query("value")))); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GinCsrfPositive(c *gin.Context) {
	if _, err := c.Cookie("auth"); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	entity := Entity{Name: c.PostForm("name")}
	database.Save(&entity)
	c.Status(http.StatusNoContent)
}

func GinCsrfControl(c *gin.Context) {
	if c.GetHeader("Authorization") == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	entity := Entity{Name: c.PostForm("name")}
	database.Save(&entity)
	c.Status(http.StatusNoContent)
}

func GinCsrfUncertain(c *gin.Context) {
	if _, err := c.Cookie("auth"); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !verifyRequestToken(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	entity := Entity{Name: c.PostForm("name")}
	database.Save(&entity)
	c.Status(http.StatusNoContent)
}

func GinRedirectPositive(c *gin.Context) {
	c.Redirect(http.StatusFound, c.Query("value"))
}

func GinRedirectControl(c *gin.Context) {
	c.Redirect(http.StatusFound, "/local")
}

func GinRedirectUncertain(c *gin.Context) {
	c.Redirect(http.StatusFound, opaqueValue(c.Query("value")))
}

func GinPolicyPositive(c *gin.Context) {
	allowed, err := policy.Enforce("subject", "object", "write")
	if err != nil || !allowed {
		c.Status(http.StatusForbidden)
	}
	c.Status(http.StatusNoContent)
}

func GinPolicyControl(c *gin.Context) {
	allowed, err := policy.Enforce("subject", "object", "write")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !allowed {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Status(http.StatusNoContent)
}

func GinPolicyUncertain(c *gin.Context) {
	if gin.Mode() == gin.ReleaseMode {
		allowed, err := policy.Enforce("subject", "object", "write")
		if err != nil || !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func GinSessionPositive(c *gin.Context) {
	if _, err := c.Cookie("auth"); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.SetCookie("auth", "credential", 600, "/", "", false, true)
	c.Status(http.StatusNoContent)
}

func GinSessionControl(c *gin.Context) {
	if _, err := c.Cookie("auth"); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.SetCookie("auth", "credential", 600, "/", "", true, true)
	c.Status(http.StatusNoContent)
}

func GinSessionUncertain(c *gin.Context) {
	if _, err := c.Cookie("auth"); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if gin.Mode() != gin.ReleaseMode {
		c.SetCookie("auth", "credential", 600, "/", "", false, true)
	}
	c.Status(http.StatusNoContent)
}

func GinCorsPositive(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GinCorsControl(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "https://trusted.example")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GinCorsUncertain(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", opaqueValue(c.Request.Header.Get("Origin")))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GinBindingPositive(c *gin.Context) {
	var entity Entity
	entity.Role = "member"
	if err := c.ShouldBind(&entity); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	database.Save(&entity)
	c.JSON(http.StatusCreated, entity)
}

func GinBindingControl(c *gin.Context) {
	var entity Entity
	entity.Role = "member"
	if err := c.ShouldBind(&entity); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	entity.Role = "member"
	database.Save(&entity)
	c.JSON(http.StatusCreated, entity)
}

func GinBindingUncertain(c *gin.Context) {
	var entity Entity
	entity.Role = "member"
	if err := c.ShouldBind(&entity); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !isAdministrator(c) {
		entity.Role = "member"
	}
	database.Save(&entity)
	c.JSON(http.StatusCreated, entity)
}
