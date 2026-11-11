package gochi

import (
	"html"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterTypedContexts(router chi.Router) {
	router.Get("/js/positive", ChiJsPositive)
	router.Get("/js/control", ChiJsControl)
	router.Get("/js/uncertain", ChiJsUncertain)
	router.Get("/js/html-encoder", ChiJsWrongHtmlEncoder)
	router.Get("/url/positive", ChiUrlPositive)
	router.Get("/url/control", ChiUrlControl)
	router.Get("/url/uncertain", ChiUrlUncertain)
	router.Get("/url/html-encoder", ChiUrlWrongHtmlEncoder)
}

func ChiJsPositive(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(w, template.JS(r.URL.Query().Get("value"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiJsControl(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(w, template.JS("const answer = 1")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiJsUncertain(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(w, template.JS(opaqueValue(r.URL.Query().Get("value")))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiJsWrongHtmlEncoder(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("script").Parse(`<script>{{.}}</script>`))
	if err := t.Execute(w, template.JS(html.EscapeString(r.URL.Query().Get("value")))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiUrlPositive(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(w, template.URL(r.URL.Query().Get("value"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiUrlControl(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(w, template.URL("https://trusted.example")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiUrlUncertain(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(w, template.URL(opaqueValue(r.URL.Query().Get("value")))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChiUrlWrongHtmlEncoder(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("link").Parse(`<a href="{{.}}">continue</a>`))
	if err := t.Execute(w, template.URL(html.EscapeString(r.URL.Query().Get("value")))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
