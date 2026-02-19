package templates

import (
	"html/template"
	"log"
	"net/http"
)

type listItem struct {
	Template *template.Template
	LayoutName string
}

var list = make(map[string]listItem)

func Parse()  {
	var Layouts *template.Template
  /* Layouts */
	log.Println("Parsing layouts...")
	parsedLayouts, err := template.ParseGlob("templates/layouts/*.html")
	if err != nil {
		panic(err)
	}
	Layouts = parsedLayouts

	type TemplateDeclaration struct {
		Layout string
		Partial string
	}

	Arg := map[string]TemplateDeclaration{
		"login": {
			Layout: "login.html",
			Partial: "templates/partial/login.html",
		},
		"foo": {
			Layout: "login.html",
			Partial: "templates/partial/about.html",
		},
	}

  /* Templates */
	for key, value := range Arg {
		clonedLayouts, err := Layouts.Clone()
		if err != nil {
			panic(err)
		}
		partialParsed, err := clonedLayouts.ParseFiles(value.Partial)
		if err != nil {
			panic(err)
		}
		list[key] = listItem{ partialParsed, value.Layout }
	}
	
	/*
	log.Println("Parsing templates...")
	clonedLayouts, err := Layouts.Clone()
	if err != nil {
		panic(err)
	}
	list["login"], err = clonedLayouts.ParseFiles("templates/partial/login.html")
	if err != nil {
		panic(err)
	}
	*/
}

func Execute(writer http.ResponseWriter, name string, data any) {
	// log.Fatal(name, list[name])
	list[name].Template.ExecuteTemplate(writer, list[name].LayoutName, data)
}
