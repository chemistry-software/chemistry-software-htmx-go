package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	IsDark bool
}

// Custom template functions
var templateFuncs = template.FuncMap{
	"iterate": func(count int) []int {
		var i []int
		for j := 0; j < count; j++ {
			i = append(i, j)
		}
		return i
	},
	"now": time.Now, // Make current time available in templates
}

// Logging middleware
func loggingMiddleware(next http.Handler, isDev bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDev {
			start := time.Now()
			log.Printf("Request: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
			log.Printf("Response: %s %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		} else {
			next.ServeHTTP(w, r) // Just serve the request without logging in production
		}
	})
}

func main() {
	devMode := flag.Bool("dev", false, "Enable development mode with verbose logging")
	flag.Parse()

	// Parse templates with custom functions
	tmpl, err := template.New("").Funcs(templateFuncs).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", loggingMiddleware(fs, *devMode)))

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		tmpl.ExecuteTemplate(w, "index.html", PageData{IsDark: isDark})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		if err := tmpl.ExecuteTemplate(w, "about.html", PageData{IsDark: isDark}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		if err := tmpl.ExecuteTemplate(w, "services.html", PageData{IsDark: isDark}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		if err := tmpl.ExecuteTemplate(w, "contact.html", PageData{IsDark: isDark}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Theme toggle endpoint
	http.HandleFunc("/toggle-theme", func(w http.ResponseWriter, r *http.Request) {
		isDark := r.URL.Query().Get("dark") == "true"
		setThemePreference(w, isDark) // Toggle the preference
		if isDark {
			w.Write([]byte(`
				<button
					class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"
					hx-get="/toggle-theme?dark=false"
					hx-swap="outerHTML"
				>
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
					</svg>
				</button>
			`))
		} else {
			w.Write([]byte(`
				<button
					class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"
					hx-get="/toggle-theme?dark=true"
					hx-swap="outerHTML"
				>
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
					</svg>
				</button>
			`))
		}
	})

	// Wrap the default handler with logging middleware
	defaultMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux() // Replace default mux to wrap the handlers
	http.Handle("/", loggingMiddleware(defaultMux, *devMode))

	logLevel := "Production"
	if *devMode {
		logLevel = "Development"
	}
	log.Printf("%s server starting on port 8080", logLevel)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Helper function to get theme preference from cookie
func getThemePreference(r *http.Request) bool {
	cookie, err := r.Cookie("theme")
	if err != nil {
		log.Println("Theme cookie not found, defaulting to light mode")
		return false // Default to light mode if no cookie
	}
	log.Println("Theme cookie found:", cookie.Value)
	return cookie.Value == "dark"
}

// Helper function to set theme preference in cookie
func setThemePreference(w http.ResponseWriter, isDark bool) {
	theme := "light"
	if isDark {
		theme = "dark"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    theme,
		Path:     "/",
		MaxAge:   3600 * 24 * 365, // 1 year
		HttpOnly: true,
	})
	log.Println("Set theme cookie:", theme)
}
