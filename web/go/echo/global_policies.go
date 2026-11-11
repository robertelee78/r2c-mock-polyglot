package echoweb

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var reflectingOriginPolicy = middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOriginFunc: func(origin string) (bool, error) { return true, nil }})

var trustedOriginPolicies = middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOrigins: []string{"https://trusted.example"}})

var opaqueOriginPolicy = middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOriginFunc: func(origin string) (bool, error) { return opaqueValue(origin) != "", nil }})

func InstallPackagePolicies(e *echo.Echo) {
	e.Group("/reflecting").Use(reflectingOriginPolicy)
	e.Group("/trusted").Use(trustedOriginPolicies)
	e.Group("/opaque").Use(opaqueOriginPolicy)
}
