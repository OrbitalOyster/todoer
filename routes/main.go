package routes

import (
	"net/http"
	"todoer/server"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	data := struct {
		Title string
	} {
		Title: "My template",
	}
	server.Templates.ExecuteTemplate(writer, "base.html", data)
}
