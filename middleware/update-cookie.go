package middleware

import (
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
)

func UpdateCookie(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		cookie := cookies.Get(req)
		if cookie != "" {
			payload, _ := jwt.GetPayload(cookie)
			tokenStr := jwt.Create(*payload)
			cookies.Set(writer, tokenStr, payload.RememberMe)
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
