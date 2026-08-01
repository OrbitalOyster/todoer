package middleware

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"
	"todoer/config"
	"todoer/server/token"
	"todoer/utils"
)

var publicURIs = []string{
	"/login",
	"/favicon.ico",
}

const redirectToURL = "/login"

func isPublicURL(URL string) bool {
	return slices.Contains(publicURIs, URL) ||
		strings.HasPrefix(URL, "/css/") ||
		strings.HasPrefix(URL, "/js/") ||
		strings.HasPrefix(URL, "/img/") ||
		strings.HasPrefix(URL, "/bootstrap/") ||
		strings.HasPrefix(URL, "/bootstrap-icons/") ||
		strings.HasPrefix(URL, "/htmx/")
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if isPublicURL(req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		/* Init token */
		token := token.Init[utils.Payload](
			req,
			&writer,
			config.CookieName,
			config.JWTSecret,
		)
		if _, err := token.Load(); err != nil {
			log.Printf("Redirect to login: %s", err)
			/* Add HTMX redirect header on HTMX requests, otherwise redirect */
			if req.Header.Get("HX-Request") == "true" {
				writer.Header().Set("HX-Redirect", redirectToURL)
			} else {
				http.Redirect(writer, req, redirectToURL, http.StatusSeeOther)
			}
		} else {
			/* Reissue the token */
			token.Save()
			/* Save token to context, pass it down the line */
			ctx := context.WithValue(req.Context(), "token", &token)
			/* Done */
			next.ServeHTTP(writer, req.WithContext(ctx))
		}
	})
}
