package api

import (
	"log"
	"net/http"
	"todoer/config"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/toasts"
)

func LoginAttempt(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	/* Credentials mock up */
	username, password, rememberMeStr :=
		req.FormValue("username"),
		req.FormValue("password"),
		req.FormValue("rememberMe")
	rememberMe := false
	if rememberMeStr == "on" {
		rememberMe = true
	}
	/* Auth mockup */
	if username == "admin" && password == "password" {
		payload := jwt.Payload{
			UserID:     username,
			RememberMe: rememberMe,
			PageSize:   config.DefaultPageSize,
			Page:       0,
			Filter:     "",
		}
		token := jwt.Set(payload)
		cookies.Set(writer, token, rememberMe)
		writer.Header().Set("HX-Redirect", "/")
		log.Printf("User %s logged in", username)
	} else {
		toasts.Info(writer, "Login failed", "Try again")
	}
}
