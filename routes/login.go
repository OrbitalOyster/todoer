package routes

import (
	"net/http"
	"todoer/templates"
)

func Login(writer http.ResponseWriter, req *http.Request) {
	templates.Execute(writer, "login", nil)
}
