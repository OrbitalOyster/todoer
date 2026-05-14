package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERR] %v\n%s", err, debug.Stack())
				http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
