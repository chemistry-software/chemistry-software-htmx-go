package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title string
}

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type Server struct {
	templates *template.Template
}

func NewServer() (*Server, error) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/about.html",
		"templates/services.html",
		"templates/contact.html",
		"templates/contact-success.html",
		"templates/projects.html",
	)
	if err != nil {
		return nil, err
	}
	log.Println("Templates loaded successfully:")
	for _, t := range tmpl.Templates() {
		log.Println(t.Name())
	}
	return &Server{templates: tmpl}, nil
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, pageName string, data PageData) {
	var err error
	if r.Header.Get("HX-Request") == "true" {
		err = s.templates.ExecuteTemplate(w, pageName, data)
	} else {
		err = s.templates.ExecuteTemplate(w, "base", data)
	}
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var form ContactForm
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// TODO: Process the contact form (e.g., send email, store in database)
		log.Printf("Received contact form: %+v", form)

		w.Header().Set("Content-Type", "text/html")
		s.templates.ExecuteTemplate(w, "contact-success", nil)
		return
	}

	s.renderTemplate(w, r, "contact", PageData{Title: "Contact"})
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		server.renderTemplate(w, r, "index", PageData{Title: "Home"})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		server.renderTemplate(w, r, "about", PageData{Title: "About"})
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		server.renderTemplate(w, r, "services", PageData{Title: "Services"})
	})

	http.HandleFunc("/contact", server.handleContact)

	// Start server
	log.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
