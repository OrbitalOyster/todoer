package routes

import (
	"log"
	"net/http"
	"todoer/server/pages"
	"todoer/server/toasts"
	"todoer/server/token"
)

/* GET */
func GetLoginPage(writer http.ResponseWriter, req *http.Request) {
	data := struct{ Title string }{"Login"}
	pages.Execute(writer, "login", data)
}

/* POST */
func LoginAttempt(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	/* Credentials mock up */
	username, password, rememberMeStr, rememberMe :=
		req.FormValue("username"),
		req.FormValue("password"),
		req.FormValue("remember-me"),
		false
	if rememberMeStr == "on" {
		rememberMe = true
	}
	/* Auth mockup */
	if username == "admin" && password == "password" {
		token.CreateFresh(username, rememberMe, writer)
		writer.Header().Set("HX-Redirect", "/")
		log.Printf("User %s logged in", username)
	} else {
		toasts.Info(writer, "Login failed", "Try again")
	}
}

func Logout(writer http.ResponseWriter, req *http.Request) {
	user := token.Get(req).UserID
	token.Clear(writer)
	writer.Header().Set("HX-Redirect", "/login")
	log.Printf("User %s logged out", user)
}
