package api

import (
	"net/http"
	"todoer/cookies"
)

func Logout(writer http.ResponseWriter, req *http.Request) {
	cookies.Clear(writer)
	http.Redirect(writer, req, "/login", http.StatusSeeOther)
}
