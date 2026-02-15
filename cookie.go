package main

import (
	"net/http"
	"time"
)

const (
	cookieName = "auth"
)

func setCookie(writer http.ResponseWriter, value string) {
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
