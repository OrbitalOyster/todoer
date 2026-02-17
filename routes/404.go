package routes

import (
	"net/http"
)

func NotFoundHandler(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("Nothing here"))
}
