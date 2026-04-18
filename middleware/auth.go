package middleware

import (
	"log"
	"net/http"
	"todoer/jwt"
	"todoer/utils"
)

func Auth(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if utils.IsPublicURL(req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		_, err := jwt.Get(req)
		if err != nil {
			log.Printf("Redirecting user to login: %s", err)
			if req.Header.Get("HX-Request") == "true" {
				writer.Header().Set("HX-Redirect", "/login")
			} else {
				http.Redirect(writer, req, "/login", http.StatusSeeOther)
			}
			return
		}
		/* All good */
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
