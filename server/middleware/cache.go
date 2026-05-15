package middleware

import (
	"net/http"
)

func Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		/* 30 days cache */
		writer.Header().Set("Cache-Control", "max-age=2592000, public")
		next.ServeHTTP(writer, req)
	})
}
