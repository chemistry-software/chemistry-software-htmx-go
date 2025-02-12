package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type PageData struct {
	IsDark bool
	Lang   string
	Trans  *message.Printer
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
	"t": func(p *message.Printer, key message.Reference, args ...interface{}) string {
		return p.Sprintf(key, args...)
	},
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

	// Localization setup
	english := language.English
	dutch := language.Dutch

	// Catalog for messages
	var mc = catalog.NewBuilder()

	// Load translations from JSON files
	loadTranslations(mc, english, "i18n/en.json")
	loadTranslations(mc, dutch, "i18n/nl.json")

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
		langTag := getLanguagePreference(r, w)
		p := message.NewPrinter(langTag, message.Catalog(mc))

		data := PageData{IsDark: isDark, Lang: langTag.String(), Trans: p}
		log.Printf("PageData: %+v", data)
		tmpl.ExecuteTemplate(w, "index.html", data)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		langTag := getLanguagePreference(r, w)
		p := message.NewPrinter(langTag, message.Catalog(mc))

		data := PageData{IsDark: isDark, Lang: langTag.String(), Trans: p}
		if err := tmpl.ExecuteTemplate(w, "about.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		langTag := getLanguagePreference(r, w)
		p := message.NewPrinter(langTag, message.Catalog(mc))
		data := PageData{IsDark: isDark, Lang: langTag.String(), Trans: p}
		log.Printf("PageData: %+v", data)
		if err := tmpl.ExecuteTemplate(w, "services.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		isDark := getThemePreference(r)
		langTag := getLanguagePreference(r, w)
		p := message.NewPrinter(langTag, message.Catalog(mc))
		data := PageData{IsDark: isDark, Lang: langTag.String(), Trans: p}
		if err := tmpl.ExecuteTemplate(w, "contact.html", data); err != nil {
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

	// Language toggle endpoint
	http.HandleFunc("/toggle-language", func(w http.ResponseWriter, r *http.Request) {
		currentLang := getLanguagePreference(r, w).String()
		var newLangTag language.Tag
		if currentLang == "en" {
			newLangTag = dutch
		} else {
			newLangTag = english
		}
		setLanguagePreference(w, newLangTag)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther) // Redirect back to the referring page
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

// Helper function to load translations from JSON file
func loadTranslations(mc *catalog.Builder, langTag language.Tag, filename string) {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading translation file %s: %v", filename, err)
	}

	var translations map[string]string
	err = json.Unmarshal(jsonFile, &translations)
	if err != nil {
		log.Fatalf("Error parsing JSON from file %s: %v", filename, err)
	}

	for key, value := range translations {
		mc.SetString(langTag, key, value)
	}
	log.Printf("Loaded %d translations for language %s from %s", len(translations), langTag, filename)
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
		Name:   "theme",
		Value:  theme,
		Path:   "/",
		MaxAge: 3600 * 24 * 365, // 1 year
		// HttpOnly: true, <- makes cookie not accessible to JS (removed for theme.js to work)
	})
	log.Println("Set theme cookie:", theme)
}

// Helper function to get language preference from cookie or Accept-Language header
func getLanguagePreference(r *http.Request, w http.ResponseWriter) language.Tag {
	cookie, err := r.Cookie("lang")
	if err == nil {
		langTag, err := language.Parse(cookie.Value)
		if err == nil {
			log.Println("Language cookie found:", langTag)
			return langTag // Return language from cookie if valid
		}
		log.Println("Invalid language cookie value:", cookie.Value)
	}

	// Fallback to Accept-Language header
	acceptLang := r.Header.Get("Accept-Language")
	userLangs, _, err := language.ParseAcceptLanguage(acceptLang)
	if err != nil {
		log.Println("Error parsing Accept-Language header:", err)
		userLangs = []language.Tag{language.English} // Default to English on error
	}

	// Check if Dutch is preferred, otherwise default to English
	preferredLang := language.English
	for _, lang := range userLangs {
		if lang.String() == "nl" {
			preferredLang = language.Dutch
			break
		}
	}

	// Set language cookie if not already set or invalid
	setLanguagePreference(w, preferredLang)
	log.Println("Setting language cookie to:", preferredLang)
	return preferredLang
}

// Helper function to set language preference in cookie
func setLanguagePreference(w http.ResponseWriter, langTag language.Tag) {
	http.SetCookie(w, &http.Cookie{
		Name:   "lang",
		Value:  langTag.String(),
		Path:   "/",
		MaxAge: 3600 * 24 * 365, // 1 year
		// HttpOnly: true, <-  Let's keep it accessible to JS if needed, adjust if security is a concern
	})
	log.Println("Set language cookie:", langTag.String())
}
