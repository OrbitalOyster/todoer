package api

import (
	"net/http"
	"todoer/server"
)

func LoginAttemptHandler(writer http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}

	username, password := req.FormValue("username"), req.FormValue("password")

	if username == "orbital" && password == "qwerty" {
		token := server.CreateToken("orbital")
		server.SetCookie(writer, token)
		writer.Header().Set("HX-Redirect", "/")
		writer.Write([]byte("Yay!"))
	} else {
		writer.Write([]byte("Try again"))
	}
}

