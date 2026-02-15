package main

import (
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	logHandler := func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s %s\n", req.Method, req.RemoteAddr, req.URL.Path)
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(logHandler)
}

