package routes

import (
	"net/http"
	"todoer/templates"
)

func Login(writer http.ResponseWriter, req *http.Request) {
	data := struct { Title string } { "Login" }
	templates.Execute(writer, "login", data)
}
