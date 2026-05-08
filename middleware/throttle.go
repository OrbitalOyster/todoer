package middleware

import (
	"net/http"
	"time"
)

func Throttle(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		if req.Header.Get("HX-Request") == "true" {
			time.Sleep(time.Second / 5)
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
