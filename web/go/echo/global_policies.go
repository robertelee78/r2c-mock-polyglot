package echoweb

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	reflectingOriginPolicy echo.MiddlewareFunc
	trustedOriginPolicies  echo.MiddlewareFunc
	opaqueOriginPolicy     echo.MiddlewareFunc
)

func init() {
	reflectingOriginPolicy = middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOriginFunc: func(origin string) (bool, error) { return true, nil }})
}

func init() {
	trustedOriginPolicies = middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOrigins: []string{"https://trusted.example"}})
}

func init() {
	opaqueOriginPolicy = middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOriginFunc: func(origin string) (bool, error) { return opaqueValue(origin) != "", nil }})
}

func InstallPackagePolicies(e *echo.Echo) {
	e.Group("/reflecting").Use(reflectingOriginPolicy)
	e.Group("/trusted").Use(trustedOriginPolicies)
	e.Group("/opaque").Use(opaqueOriginPolicy)
}
