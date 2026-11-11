package ginweb

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SessionRouter() *gin.Engine {
	router := gin.New()
	router.Use(sessions.Sessions("auth", cookie.NewStore([]byte(os.Getenv("SESSION_KEY")))))
	router.Group("/plain").Use(GinSessionOptionsPositive)
	router.Group("/secure").Use(GinSessionOptionsControl)
	router.POST("/entities", func(c *gin.Context) {
		GinCsrfTokenControl(c, database)
	})
	return router
}

func GinCsrfTokenControl(c *gin.Context, database *gorm.DB) {
	session := sessions.Default(c)
	if c.GetHeader("X-CSRF-Token") != session.Get("csrf") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	entity := Entity{Name: c.PostForm("name")}
	database.Save(&entity)
	c.Status(http.StatusNoContent)
}

func GinSessionOptionsPositive(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{Secure: false, HttpOnly: true})
	c.Next()
}

func GinSessionOptionsControl(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{Secure: true, HttpOnly: true})
	c.Next()
}

func GinCorsPositiveRouter() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{AllowCredentials: true, AllowOriginFunc: func(origin string) bool { return true }}))
	return router
}

func GinCorsControlRouter() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: []string{"https://trusted.example"}}))
	return router
}

func GinCorsUncertainRouter() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{AllowCredentials: true, AllowOriginFunc: func(origin string) bool { return opaqueValue(origin) != "" }}))
	return router
}
