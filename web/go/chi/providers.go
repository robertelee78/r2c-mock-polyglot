package gochi

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

var csrfAuthKey = []byte(os.Getenv("CSRF_AUTH_KEY"))

func CsrfRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(csrf.Protect(csrfAuthKey))
	router.Post("/entities", ChiInstalledCsrfControl)
	router.Get("/session/options/positive", ChiSessionOptionsPositive)
	router.Get("/session/options/control", ChiSessionOptionsControl)
	return router
}

func requestDatabase(r *http.Request) *gorm.DB {
	return database.WithContext(r.Context())
}

func ChiInstalledCsrfControl(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	database := requestDatabase(r)
	entity := Entity{Name: r.FormValue("name")}
	database.Save(&entity)
	w.WriteHeader(http.StatusNoContent)
}

func ChiSessionOptionsPositive(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{Secure: false, HttpOnly: true}
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiSessionOptionsControl(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{Secure: true, HttpOnly: true}
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiCorsPositiveRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{AllowCredentials: true, AllowOriginFunc: func(request *http.Request, origin string) bool { return true }}))
	return router
}

func ChiCorsControlRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{AllowCredentials: true, AllowedOrigins: []string{"https://trusted.example"}}))
	return router
}

func ChiCorsUncertainRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{AllowCredentials: true, AllowOriginFunc: func(request *http.Request, origin string) bool { return opaqueValue(origin) != "" }}))
	return router
}
