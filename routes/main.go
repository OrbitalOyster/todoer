package routes

import (
	"net/http"
	"todoer/templates"
	"todoer/jwt"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	claims := jwt.Get(req)
	data := struct {
		Title string
		Username string
	} {
		Title: "todoer",
		Username: claims.UserID,
	}
	templates.Execute(writer, "main", data)
}
