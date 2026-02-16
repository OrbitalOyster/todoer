package middlware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s %s\n", req.Method, req.RemoteAddr, req.URL.Path)
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
