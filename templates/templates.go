package templates

import (
	"html/template"
	"log"
	"net/http"
)

type TemplateDeclaration struct {
	Layout  string
	Partial string
}

var list = make(map[string]struct {
	*template.Template
	layoutName string
})

func Parse(templateMap map[string]TemplateDeclaration) {
	/* Layouts */
	log.Println("Parsing layouts...")
	layouts, err := template.ParseGlob("templates/layouts/*.html")
	if err != nil {
		panic(err)
	}
	/* Templates */
	log.Println("Parsing templates...")
	for key, value := range templateMap {
		clonedLayouts, err := layouts.Clone()
		if err != nil {
			panic(err)
		}
		partialParsed, err := clonedLayouts.ParseFiles("templates/partial/" + value.Partial)
		if err != nil {
			panic(err)
		}
		list[key] = struct {
			*template.Template
			layoutName string
		}{
			partialParsed,
			value.Layout,
		}
		log.Printf("%s: %s - %s", key, value.Layout, value.Partial)
	}
}

func Execute(writer http.ResponseWriter, name string, data any) {
	list[name].ExecuteTemplate(writer, list[name].layoutName, data)
}
