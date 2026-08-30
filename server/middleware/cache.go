package middleware

import (
	"fmt"
	"net/http"
)

/* 30 days cache */
const maxAgeSeconds = 60 * 60 * 24 * 30

func Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set(
			"Cache-Control",
			fmt.Sprintf("max-age=%d, public", maxAgeSeconds),
		)
		next.ServeHTTP(writer, req)
	})
}
