package gochi

import (
	"encoding/json"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Entity struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var database *gorm.DB

var policy *casbin.Enforcer

func RegisterCases(router chi.Router) {
	router.Get("/html/positive", ChiHtmlPositive)
	router.Get("/html/control", ChiHtmlControl)
	router.Get("/html/uncertain", ChiHtmlUncertain)
	router.Post("/csrf/positive", ChiCsrfPositive)
	router.Post("/csrf/control", ChiCsrfControl)
	router.Post("/csrf/uncertain", ChiCsrfUncertain)
	router.Get("/redirect/positive", ChiRedirectPositive)
	router.Get("/redirect/control", ChiRedirectControl)
	router.Get("/redirect/uncertain", ChiRedirectUncertain)
	router.Post("/policy/positive", ChiPolicyPositive)
	router.Post("/policy/control", ChiPolicyControl)
	router.Post("/policy/uncertain", ChiPolicyUncertain)
	router.Post("/session/positive", ChiSessionPositive)
	router.Post("/session/control", ChiSessionControl)
	router.Post("/session/uncertain", ChiSessionUncertain)
	router.Get("/cors/positive", ChiCorsPositive)
	router.Get("/cors/control", ChiCorsControl)
	router.Get("/cors/uncertain", ChiCorsUncertain)
	router.Post("/binding/positive", ChiBindingPositive)
	router.Post("/binding/control", ChiBindingControl)
	router.Post("/binding/uncertain", ChiBindingUncertain)
}

func ChiHtmlPositive(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	if err := t.Execute(w, template.HTML(r.URL.Query().Get("value"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiHtmlControl(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	if err := t.Execute(w, template.HTML(html.EscapeString(r.URL.Query().Get("value")))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiHtmlUncertain(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("page").Parse(`<main>{{.}}</main>`))
	if err := t.Execute(w, template.HTML(opaqueValue(r.URL.Query().Get("value")))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiCsrfPositive(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	entity := Entity{Name: r.FormValue("name")}
	database.Save(&entity)
	w.WriteHeader(http.StatusNoContent)
}

func ChiCsrfControl(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	entity := Entity{Name: r.FormValue("name")}
	database.Save(&entity)
	w.WriteHeader(http.StatusNoContent)
}

func ChiCsrfUncertain(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if !verifyRequestToken(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	entity := Entity{Name: r.FormValue("name")}
	database.Save(&entity)
	w.WriteHeader(http.StatusNoContent)
}

func ChiRedirectPositive(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Query().Get("value"), http.StatusFound)
}

func ChiRedirectControl(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/local", http.StatusFound)
}

func ChiRedirectUncertain(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, opaqueValue(r.URL.Query().Get("value")), http.StatusFound)
}

func ChiPolicyPositive(w http.ResponseWriter, r *http.Request) {
	allowed, err := policy.Enforce("subject", "object", "write")
	if err != nil || !allowed {
		log.Printf("policy rejected %s %s", r.Method, r.URL.Path)
	}
	w.WriteHeader(http.StatusNoContent)
}

func ChiPolicyControl(w http.ResponseWriter, r *http.Request) {
	allowed, err := policy.Enforce("subject", "object", "write")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ChiPolicyUncertain(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("POLICY_ENFORCEMENT") != "disabled" {
		allowed, err := policy.Enforce("subject", "object", "write")
		if err != nil || !allowed {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func ChiSessionPositive(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "auth", Value: "credential", Secure: false, HttpOnly: true})
	w.WriteHeader(http.StatusNoContent)
}

func ChiSessionControl(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "auth", Value: "credential", Secure: true, HttpOnly: true})
	w.WriteHeader(http.StatusNoContent)
}

func ChiSessionUncertain(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if os.Getenv("APP_ENV") == "development" {
		http.SetCookie(w, &http.Cookie{Name: "auth", Value: "credential", Secure: false, HttpOnly: true})
	}
	w.WriteHeader(http.StatusNoContent)
}

func ChiCorsPositive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusNoContent)
}

func ChiCorsControl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://trusted.example")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusNoContent)
}

func ChiCorsUncertain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", opaqueValue(r.Header.Get("Origin")))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusNoContent)
}

func ChiBindingPositive(w http.ResponseWriter, r *http.Request) {
	var entity Entity
	entity.Role = "member"
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database.Save(&entity)
	w.WriteHeader(http.StatusCreated)
}

func ChiBindingControl(w http.ResponseWriter, r *http.Request) {
	var entity Entity
	entity.Role = "member"
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entity.Role = "member"
	database.Save(&entity)
	w.WriteHeader(http.StatusCreated)
}

func ChiBindingUncertain(w http.ResponseWriter, r *http.Request) {
	var entity Entity
	entity.Role = "member"
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isAdministrator(r) {
		entity.Role = "member"
	}
	database.Save(&entity)
	w.WriteHeader(http.StatusCreated)
}
