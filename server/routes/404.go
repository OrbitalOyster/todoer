package routes

import (
	"net/http"
)

func NotFoundHandler(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	_, err := writer.Write([]byte("Nothing here"))
	if err != nil {
		panic(err)
	}
}
