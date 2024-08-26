package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	parsedTemplate, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", tmpl),
	)
	if err != nil {
		log.Printf("Error parsing template files: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = parsedTemplate.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index.html", map[string]string{"Title": "Home"})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about.html", map[string]string{"Title": "About Us"})
	})

	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "projects.html", map[string]string{"Title": "Our Projects"})
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "contact.html", map[string]string{"Title": "Contact Us"})
	})

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
