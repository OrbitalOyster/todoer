package routes

import (
	"net/http"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/html/index.html")
}
