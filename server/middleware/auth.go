package middleware

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"todoer/config"
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

		funkyToken := token.CreateFunkyToken[token.Payload](
			req,
			&writer,
			config.CookieName,
			config.JWTSecret,
			config.CookieShortLifetime,
		)

		// if err := funky.Load(); err != nil {
		// 	log.Printf("ERROR: %v", err)
		// } else {
		// 	log.Printf("RESULT: %v\n", funky.GetPayload())
		// }

		/* Protected routes - check credentials */
		// GetPayloadSafe := func() *token.Payload {
		// 	/* On fail - redirect to login */
		// 	defer func() {
		// 		if r := recover(); r != nil {
		// 			log.Printf("Redirect to login: %s", r)
		// 			/* Add HTMX redirect header on HTMX requests, otherwise redirect */
		// 			if req.Header.Get("HX-Request") == "true" {
		// 				writer.Header().Set("HX-Redirect", "/login")
		// 			} else {
		// 				http.Redirect(writer, req, "/login", http.StatusSeeOther)
		// 			}
		// 		}
		// 	}()
		// 	return token.Get(req)
		// }

		if err := funkyToken.Load(); err != nil {
			log.Printf("Redirect to login: %s", err)
			/* Add HTMX redirect header on HTMX requests, otherwise redirect */
			if req.Header.Get("HX-Request") == "true" {
				writer.Header().Set("HX-Redirect", "/login")
			} else {
				http.Redirect(writer, req, "/login", http.StatusSeeOther)
			}
		} else {
			/* Reissue the token */
			funkyToken.Save()
			/* Done */
			next.ServeHTTP(writer, req)
		}

		// if payload := GetPayloadSafe(); payload != nil {
		// 	/* Reissue the token */
		// 	token.Create(*payload, writer)
		// 	/* Done */
		// 	next.ServeHTTP(writer, req)
		// }
	})
}
