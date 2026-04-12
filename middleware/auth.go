package middleware

import (
	"net/http"
	"slices"
	"strings"
	"todoer/jwt"
)

var publicURIs = []string{
	"/login",
	"/favicon.ico",
	"/css/reset.css",
	"/css/style.css",
	"/js/script.js",
}

func Auth(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if slices.Contains(publicURIs, req.URL.Path) || strings.HasPrefix(req.URL.Path, "/vendor/") {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		_, err := jwt.Get(req)
		if err != nil {
			http.Redirect(writer, req, "/login", http.StatusSeeOther)
			return
		}
		/* All good */
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
