package middleware

import (
	"log"
	"net/http"
	"todoer/jwt"
	"todoer/utils"
)

func Token(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		/* On fail - redirect to login */
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Redirect to login: %s", r)
				/* Add HTMX redirect header on HTMX requests, otherwise redirect */
				if req.Header.Get("HX-Request") == "true" {
					writer.Header().Set("HX-Redirect", "/login")
				} else {
					http.Redirect(writer, req, "/login", http.StatusSeeOther)
				}
			}
		}()
		/* Public routes */
		if utils.IsPublicURL(req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		payload := jwt.Get(req)
		/* Reissue the token */
		jwt.Create(*payload, writer)
		/* Done */
		next.ServeHTTP(writer, req)
	})
}
