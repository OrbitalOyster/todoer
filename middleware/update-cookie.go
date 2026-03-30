package middlware

import (
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
)

func UpdateCookie(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Protected routes - check credentials */
		cookie := cookies.Get(req)
		/* No cookie */
		if cookie != "" {
			/* JWT is already validated */
			claims, _ := jwt.Validate(cookie)
			userId := claims.UserID
			rememberMe := claims.RememberMe
			token := jwt.Create(userId, rememberMe)
			cookies.Set(writer, token, rememberMe)
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
