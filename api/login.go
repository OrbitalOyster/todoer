package api

import (
	"log"
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/templates"
)

const loginFailedMsg = `{
	"toast": {
		"type": "danger",
		"title": "Login failed",
		"msg": "Try again"
	}
}`

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
		token := jwt.Create(username, rememberMe)
		cookies.Set(writer, token, rememberMe)
		writer.Header().Set("HX-Redirect", "/")
		log.Printf("User %s logged in", username)
	} else {
		// writer.Header().Set("HX-Trigger", loginFailedMsg)
		writer.Header().Set("HX-Trigger-After-Settle", `{"createToast":{"target" : ".toast-container"}}`)
		// writer.WriteHeader(http.StatusNoContent)
		templates.ExecutePartial(writer, "foo", nil)	
	}
}
