package app

import (
	"html/template"
	"log"
	"net/http"
)

var (
	templates      *template.Template
	CurrentVersion string
)

// InitTemplates parses the templates and panics if it can't
func InitTemplates(path string) {
	templates = template.Must(template.ParseGlob(path))
	log.Println("Initialized templates")
}

// RenderTemplate writes a template to a Response
// Panics if the templates haven't been initialized.
func RenderTemplate(w http.ResponseWriter, template string, data interface{}) {
	if templates == nil {
		panic("Templates unitialized")
	}
	err := templates.ExecuteTemplate(w, template+".html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
	}
}
