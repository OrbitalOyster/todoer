package middlware

import (
	"net/http"
	"slices"
	"todoer/server"
)

var	publicURIs = []string {
	"/login",
	"/favicon.ico",
	"/css/reset.css",
}

func Auth(next http.Handler) http.Handler  {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if slices.Contains(publicURIs, req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		cookie := server.GetCookie(req)
		/* No cookie */
		if cookie == "" {
			http.Redirect(writer, req, "/login", http.StatusSeeOther)
			return
		}
		_, err := server.ValidateToken(cookie)
		if err != nil {
			http.Redirect(writer, req, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
