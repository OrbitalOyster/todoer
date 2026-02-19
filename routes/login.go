package routes

import (
	"net/http"
	"todoer/templates"
)

func Login(writer http.ResponseWriter, req *http.Request) {
	// http.ServeFile(writer, req, "static/html/login.html")
	templates.Execute(writer, "login", nil)
}
