package gochi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

var reflectingCorsHandler = cors.Handler(cors.Options{AllowCredentials: true, AllowOriginFunc: func(request *http.Request, origin string) bool { return true }})

var trustedOriginHandler = cors.Handler(cors.Options{AllowCredentials: true, AllowedOrigins: []string{"https://trusted.example"}})

var opaqueOriginHandler = cors.Handler(cors.Options{AllowCredentials: true, AllowOriginFunc: func(request *http.Request, origin string) bool { return opaqueValue(origin) != "" }})

func InstallPackagePolicies(router chi.Router, api http.Handler) {
	router.With(reflectingCorsHandler).Mount("/reflecting", api)
	router.With(trustedOriginHandler).Mount("/trusted", api)
	router.With(opaqueOriginHandler).Mount("/opaque", api)
}
