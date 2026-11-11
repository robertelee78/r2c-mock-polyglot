package gochi

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	texttemplate "text/template"

	"github.com/go-chi/chi/v5"
)

type sessionConfig struct {
	Secure bool
}

func RegisterSupporting(router chi.Router) {
	router.Get("/templates/raw", ChiRawTemplate)
	router.Get("/templates/escaped", ChiEscapedTemplate)
	router.Get("/templates/text", ChiTextTemplate)
	router.Post("/session/configured", ChiConfiguredCookie)
}

func ChiRawTemplate(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("page.gotmpl", "partial.gotmpl"))
	if err := t.Execute(w, template.HTML(r.URL.Query().Get("value"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiEscapedTemplate(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("page.gotmpl", "partial.gotmpl"))
	if err := t.Execute(w, r.URL.Query().Get("value")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiTextTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := texttemplate.Must(texttemplate.ParseFiles("page.gotmpl", "partial.gotmpl"))
	if err := t.Execute(w, r.URL.Query().Get("value")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiConfiguredCookie(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	raw, err := os.ReadFile("session.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var config sessionConfig
	if err := json.Unmarshal(raw, &config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "auth", Value: "credential", Secure: config.Secure, HttpOnly: true})
	w.WriteHeader(http.StatusNoContent)
}

func opaqueValue(value string) string {
	return strings.TrimSpace(value)
}

func verifyRequestToken(r *http.Request) bool {
	return r.Header.Get("X-CSRF-Token") != ""
}

func isAdministrator(r *http.Request) bool {
	return r.Header.Get("X-Role") == "administrator"
}
