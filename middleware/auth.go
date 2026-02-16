package middlware

import (
	"net/http"
	"slices"
	"todoer/cookies"
)

var publicURIs = []string{
	"/login",
	"/favicon.ico",
	"/css/reset.css",
}

func Auth(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if slices.Contains(publicURIs, req.URL.Path) {
			next.ServeHTTP(writer, req)
			return
		}
		/* Protected routes - check credentials */
		cookie := cookies.Get(req)
		/* No cookie */
		if cookie == "" {
			http.Redirect(writer, req, "/login", http.StatusSeeOther)
			return
		}
		/* Invalid or expired token */
		_, err := cookies.ValidateToken(cookie)
		if err != nil {
			http.Redirect(writer, req, "/login", http.StatusSeeOther)
			return
		}
		/* All good */
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
