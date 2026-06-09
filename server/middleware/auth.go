package middleware

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"todoer/server/token"
)

var publicURIs = []string{
	"/login",
	"/favicon.ico",
}

func isPublicURL(URL string) bool {
	return slices.Contains(publicURIs, URL) ||
		strings.HasPrefix(URL, "/css/") ||
		strings.HasPrefix(URL, "/js/") ||
		strings.HasPrefix(URL, "/img/") ||
		strings.HasPrefix(URL, "/vendor/")
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if isPublicURL(req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		GetPayloadSafe := func() *token.Payload {
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
			return token.Get(req)
		}
		if payload := GetPayloadSafe(); payload != nil {
			/* Reissue the token */
			token.Create(*payload, writer)
			/* Done */
			next.ServeHTTP(writer, req)
		}
	})
}
