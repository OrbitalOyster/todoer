package middlware

import (
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
)

func UpdateCookie(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		cookie := cookies.Get(req)
		if cookie != "" {
			/* JWT is already validated */
			claims, _ := jwt.Validate(cookie)
			userId := claims.UserID
			rememberMe := claims.RememberMe
			pageSize := claims.PageSize
			token := jwt.Create(userId, rememberMe, pageSize)
			cookies.Set(writer, token, rememberMe)
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
