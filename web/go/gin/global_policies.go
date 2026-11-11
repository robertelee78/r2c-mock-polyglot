package ginweb

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var reflectingCorsHandler = cors.New(cors.Config{AllowCredentials: true, AllowOriginFunc: func(origin string) bool { return true }})

var trustedOriginHandler = cors.New(cors.Config{AllowCredentials: true, AllowOrigins: []string{"https://trusted.example"}})

var opaqueOriginHandler = cors.New(cors.Config{AllowCredentials: true, AllowOriginFunc: func(origin string) bool { return opaqueValue(origin) != "" }})

func InstallPackagePolicies(router *gin.Engine) {
	router.Group("/reflecting").Use(reflectingCorsHandler)
	router.Group("/trusted").Use(trustedOriginHandler)
	router.Group("/opaque").Use(opaqueOriginHandler)
}
