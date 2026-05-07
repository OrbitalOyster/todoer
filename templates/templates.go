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

const timeFormat = "2.01.2006 15:04:05"

var funcMap = template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format(timeFormat)
	},
	"greet": func(name string) string {
		return "Hello, " + name + "!"
	},
}

var (
	layoutsFolder = filepath.Join("templates", "layouts")
	pagesFolder   = filepath.Join("templates", "pages")
	partialsGlob  = filepath.Join("templates", "partials", "*.html")
	modalsGlob    = filepath.Join("templates", "partials", "modals", "*.html")
	pages         map[string]*template.Template
	partials      *template.Template
)

func init() {
	pages = make(map[string]*template.Template)
	partials = template.New("partials").Funcs(funcMap)
	template.Must(partials.ParseGlob(partialsGlob))
	template.Must(partials.ParseGlob(modalsGlob))
}

func AddPage(page string, layout string) {
	layoutFullFilename := filepath.Join(layoutsFolder, layout+".html")
	pageFullFilename := filepath.Join(pagesFolder, page+".html")
	/* Create new template, add custom functions */
	newPage := template.New(layout + ".html").Funcs(funcMap)
	template.Must(newPage.ParseFiles(layoutFullFilename, pageFullFilename))
	/* Mix in partials */
	template.Must(newPage.ParseGlob(partialsGlob))
	template.Must(newPage.ParseGlob(modalsGlob))
	/* Done */
	pages[page] = newPage
	log.Printf("Added page \"%s/%s\" \n", page, layout)
}

func ExecutePage(writer http.ResponseWriter, name string, data any) {
	/* Check if page exists */
	if page, found := pages[name]; found {
		err := page.Execute(writer, data)
		if err != nil {
			panic(err)
		}
		writer.Header().Set("Content-Type", "text/html")
	} else {
		panic("Invalid page: " + name)
	}
}

func ExecutePartial(writer http.ResponseWriter, name string, data any) {
	err := partials.ExecuteTemplate(writer, name, data)
	if err != nil {
		panic(err)
	}
	writer.Header().Set("Content-Type", "text/html")
}
