package ginweb

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	reflectingCorsHandler gin.HandlerFunc
	trustedOriginHandler  gin.HandlerFunc
	opaqueOriginHandler   gin.HandlerFunc
)

func init() {
	reflectingCorsHandler = cors.New(cors.Config{AllowCredentials: true, AllowOriginFunc: func(origin string) bool { return true }})
}

func init() {
	trustedOriginHandler = cors.New(cors.Config{AllowCredentials: true, AllowOrigins: []string{"https://trusted.example"}})
}

func init() {
	opaqueOriginHandler = cors.New(cors.Config{AllowCredentials: true, AllowOriginFunc: func(origin string) bool { return opaqueValue(origin) != "" }})
}

func InstallPackagePolicies(router *gin.Engine) {
	router.Group("/reflecting").Use(reflectingCorsHandler)
	router.Group("/trusted").Use(trustedOriginHandler)
	router.Group("/opaque").Use(opaqueOriginHandler)
}
