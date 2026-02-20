package routes

import (
	"net/http"
	// "todoer/templates"
	"html/template"
)

func Login(writer http.ResponseWriter, req *http.Request) {
	// templates.Execute(writer, "login", nil)
	base, err := template.ParseFiles("templates/base_empty.html")
	if err != nil {
		panic(err)
	}
	data := struct { Username string } { Username: "admin" }
	// data := struct { Bar string } { Bar: "bar" }
	// base.ExecuteTemplate(writer, "base_empty", nil)
	err = base.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}
