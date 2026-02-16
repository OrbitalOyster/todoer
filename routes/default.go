package routes

import (
	"net/http"
)

func DefaultHandler(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(writer, req, "static/html/index.html")
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Nothing here"))
	}
}
