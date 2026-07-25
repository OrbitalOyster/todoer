package token

import (
	"net/http"
	"time"
)

func setCookie(name string, value string, expires time.Time, writer http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	/* Avoid double headers */
	if len(writer.Header()["Set-Cookie"]) != 0 {
		writer.Header().Del("Set-Cookie")
	}
	http.SetCookie(writer, &cookie)
}

func getCookie(name string, req *http.Request) string {
	cookie, err := req.Cookie(name)
	if err != nil {
		/* No cookie, return empty string */
		if err == http.ErrNoCookie {
			return ""
		}
	}
	return cookie.Value
}

func clearCookie(name string, writer http.ResponseWriter) {
	writer.Header().Del("Set-Cookie")
	emptyCookie := http.Cookie{
		Name:    name,
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(writer, &emptyCookie)
}
