package middleware

import (
	"net/http"
	"strings"
)

/* Set content-type to "text/html" on every htmx request */
func HTMXHeaders(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/htmx/") {
			writer.Header().Set("Content-Type", "text/html")
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
