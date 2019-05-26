package pages

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func LoadPages(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func RenderPage(w http.ResponseWriter, template string, data interface{}) error {
	return templates.ExecuteTemplate(w, template, data)
}