package routes

import (
	"net/http"
	"todoer/jwt"
	"todoer/server"
)

func Foo(writer http.ResponseWriter, req *http.Request) {
	claims := jwt.Get(req)
	data := struct {
		Title string
		Username string
		Year int
	} {
		Title: "Foo",
		Username: claims.UserID,
		Year: 2026,
	}
	server.Templates.ExecuteTemplate(writer, "base.html", data)
}
