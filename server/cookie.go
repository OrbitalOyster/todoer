package server

import (
	"net/http"
	"time"
)

const (
	cookieName = "jwt"
)

func GetCookie(req *http.Request) string  {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		/* No cookie, return empry string */
		if err == http.ErrNoCookie {
			return "" 
		}
	}
	return cookie.Value
}

func SetCookie(writer http.ResponseWriter, value string) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(jwtLifetimeMinutes) * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(writer, &cookie)
}
