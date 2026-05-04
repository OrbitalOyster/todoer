package middleware

import (
	"log"
	"net/http"
	"todoer/jwt"
	"todoer/utils"
)

func Token(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if utils.IsPublicURL(req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		payload, err := jwt.Get(req)
		if err != nil {
			log.Printf("Redirecting user to login: %s", err)
			/* Add HTMX redirect header on HTMX requests, otherwise redirect */
			if req.Header.Get("HX-Request") == "true" {
				writer.Header().Set("HX-Redirect", "/login")
			} else {
				http.Redirect(writer, req, "/login", http.StatusSeeOther)
			}
			return
		}
		/* Update token */
		jwt.Create(*payload, writer)
		/* Done */
		next.ServeHTTP(writer, req)
	})
}
