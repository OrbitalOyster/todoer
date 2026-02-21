package templates

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type TemplateDeclaration struct {
	Layout  string
	Partial string
}

var layoutsGlob = filepath.Join("templates", "layouts", "*.html")
var pagesFolder = filepath.Join("templates", "pages")
var layouts *template.Template
var parsed = make(map[string]	*template.Template)

func init()  {
	/* Layouts */
	log.Printf("Parsing layouts (%s)\n", layoutsGlob)
	var err error
	layouts, err = template.ParseGlob(layoutsGlob)
	if err != nil {
		panic(err)
	}
}

func Add(name string, layout string, page string)  {
	/* Get layout */
	clonedLayouts, err := layouts.Clone()
	if err != nil {
		panic(err)
	}
	/* Get page */
	filename := filepath.Join(pagesFolder, page)
	newPage, err := clonedLayouts.ParseFiles(filename)
	if err != nil {
		panic(err)
	}
	parsed[name] = newPage
	log.Printf("Added page %s\n", name)
}

func Execute(writer http.ResponseWriter, name string, data any) {
	err := parsed[name].Execute(writer, data)
	if err != nil {
		panic(err)
	}
}
