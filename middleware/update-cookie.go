package middleware

import (
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/utils"
)

func UpdateCookie(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Skip public URLs */
		if !utils.IsPublicURL(req.URL.Path) {
			cookie := cookies.Get(req)
			if cookie != "" {
				payload, _ := jwt.GetPayload(cookie)
				tokenStr := jwt.Create(*payload)
				cookies.Set(writer, tokenStr, payload.RememberMe)
			}
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
