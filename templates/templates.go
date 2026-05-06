package templates

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type TemplateDeclaration struct {
	Layout  string
	Partial string
}

var (
	layoutsGlob  = filepath.Join("templates", "layouts", "*.html")
	partialsGlob = filepath.Join("templates", "partials", "*.html")
	modalsGlob   = filepath.Join("templates", "partials", "modals", "*.html")
	pagesFolder  = filepath.Join("templates", "pages")
	layouts      *template.Template
	partials     *template.Template
	parsed       = make(map[string]*template.Template)
)

const timeFormat = "2.01.2006 15:04:05"

var funcMap = template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format(timeFormat)
	},
	"greet": func(name string) string {
		return "Hello, " + name + "!"
	},
}

func init() {
	/* Layouts */
	layouts = template.Must(template.ParseGlob(layoutsGlob)).Funcs(funcMap)
	/* Partials */
	partials = template.Must(layouts.ParseGlob(partialsGlob))
	template.Must(partials.ParseGlob(modalsGlob))
}

func Add(name string, layout string, page string) {
	/* Clone layout */
	clonedLayouts, err := layouts.Clone()
	if err != nil {
		panic(err)
	}
	/* Parse page */
	filename := filepath.Join(pagesFolder, page)
	newPage, err := clonedLayouts.ParseFiles(filename)
	if err != nil {
		panic(err)
	}
	parsed[name] = newPage
	log.Printf("Added page \"%s\"\n", name)
}

func ExecutePage(writer http.ResponseWriter, name string, data any) {
	/* Check if page exists */
	if page, found := parsed[name]; found {
		writer.Header().Set("Content-Type", "text/html")
		err := page.Execute(writer, data)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Invalid page: " + name)
	}
}

func ExecutePartial(writer http.ResponseWriter, name string, data any) {
	writer.Header().Set("Content-Type", "text/html")
	err := partials.ExecuteTemplate(writer, name, data)
	if err != nil {
		panic(err)
	}
}
