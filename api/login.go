package api

import (
	"log"
	"net/http"
	"todoer/server"
)

func LoginAttemptHandler(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	/* Credentials mock up */
	username, password := req.FormValue("username"), req.FormValue("password")
	if username == "orbital" && password == "qwerty" {
		token := server.CreateToken(username)
		server.SetCookie(writer, token)
		writer.Header().Set("HX-Redirect", "/")
		log.Printf("User %s logged in", username)
	} else {
		writer.Write([]byte("Try again"))
	}
}
