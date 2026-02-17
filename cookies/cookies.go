package cookies

import (
	"net/http"
	"time"
	"todoer/config"
)

func Set(writer http.ResponseWriter, value string) {
	cookie := http.Cookie{
		Name:     config.CookieName,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(config.CookieLifetime) * time.Second),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(writer, &cookie)
}

func Get(req *http.Request) string {
	cookie, err := req.Cookie(config.CookieName)
	if err != nil {
		/* No cookie, return empty string */
		if err == http.ErrNoCookie {
			return ""
		}
	}
	return cookie.Value
}
