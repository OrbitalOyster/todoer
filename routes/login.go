package routes

import (
	"net/http"
)

func Login(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/html/login.html")
}
