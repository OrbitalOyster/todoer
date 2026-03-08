package api

import (
	"log"
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/templates"
	"time"
)

const loginFailedMsg = `{
	"toast": {
		"type": "danger",
		"title": "Login failed",
		"msg": "Try again"
	}
}`

func setToastHeaders(writer http.ResponseWriter) {
	writer.Header().Set("HX-Trigger-After-Settle", "toast")
	writer.Header().Set("HX-Retarget", ".toast-container")
	writer.Header().Set("HX-Reswap", "beforeend")
}

func WarningToast(writer http.ResponseWriter, title string, msg string) {
	setToastHeaders(writer)
	options := map[string]string {
		"BorderClass": "border-warning-subtle",
		"Autohide": "false",
		"HeaderColor": "bg-warning-subtle",
		"IconColor": "text-warning",
		"IconClass": "bi-exclamation-triangle-fill",
		"Title": title,
		"ProgressBarClass": "d-none",
		"Content": msg,
	}
	options["Time"] = time.Now().Format("2.01.2006 15:04:05")
	templates.ExecutePartial(writer, "toast", options)	
}

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
		WarningToast(writer, "Login failed", "Try again")
		/*
		writer.Header().Set("HX-Trigger-After-Settle", "toast")
		writer.Header().Set("HX-Retarget", ".toast-container")
		writer.Header().Set("HX-Reswap", "beforeend")
		templates.ExecutePartial(writer, "foo", nil)	
		*/
	}
}
