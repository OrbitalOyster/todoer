package cookies

import (
	"net/http"
	"time"
	"todoer/config"
)

func Set(writer http.ResponseWriter, value string, longLifetime bool) {
	expires := time.Now()
	if longLifetime {
		expires = expires.Add(time.Duration(config.CookieLifetime) * time.Second)
	} else {
		expires = expires.Add(time.Hour)
	}
	cookie := http.Cookie{
		Name:     config.CookieName,
		Value:    value,
		Expires:  expires,
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

func Clear(writer http.ResponseWriter)  {
	emptyCookie := http.Cookie{
		Name:     config.CookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(writer, &emptyCookie)
}
